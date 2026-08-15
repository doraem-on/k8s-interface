package k8sinterface

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// writeTestKubeconfig writes a kubeconfig with two contexts, each pointing at a
// distinct (unreachable) server, and points KUBECONFIG at it for the duration of
// the test. No cluster is dialed by SetClusterContextName/GetK8sConfig - this is
// about config resolution, not connectivity - so unreachable hosts are enough.
func writeTestKubeconfig(t *testing.T) (contextA, contextB string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "kubeconfig")
	const kubeconfigTemplate = `apiVersion: v1
kind: Config
clusters:
- name: cluster-a
  cluster:
    server: https://a.example.invalid:6443
- name: cluster-b
  cluster:
    server: https://b.example.invalid:6443
contexts:
- name: context-a
  context:
    cluster: cluster-a
    namespace: namespace-a
- name: context-b
  context:
    cluster: cluster-b
    namespace: namespace-b
current-context: context-a
`
	require.NoError(t, os.WriteFile(path, []byte(kubeconfigTemplate), 0o600))

	oldKubeconfig, hadKubeconfig := os.LookupEnv("KUBECONFIG")
	require.NoError(t, os.Setenv("KUBECONFIG", path))
	t.Cleanup(func() {
		if hadKubeconfig {
			_ = os.Setenv("KUBECONFIG", oldKubeconfig)
		} else {
			_ = os.Unsetenv("KUBECONFIG")
		}
	})

	return "context-a", "context-b"
}

// resetClientState clears the client-config globals SetClusterContextName's
// invalidation targets, so this test starts and ends independent of whatever ran
// before/after it in the same package (see tearDown, which does not clear K8SConfig
// since no prior test needed to).
func resetClientState(t *testing.T) {
	t.Helper()
	K8SConfig = nil
	clientConfigAPI = nil
	clusterContextName = ""
	connectedToCluster = true
}

// A second call to SetClusterContextName for a different context must invalidate
// the cached K8SConfig/clientConfigAPI, or GetK8sConfig keeps returning the first
// context's server under the second context's name (#158): GetContextName() would
// report "context-b" while GetK8sConfig().Host is still cluster-a's server, since
// LoadK8sConfig/GetConfig only (re)load when their own cache is nil.
func TestSetClusterContextName_InvalidatesConfigOnContextChange(t *testing.T) {
	resetClientState(t)
	t.Cleanup(func() { resetClientState(t) })

	contextA, contextB := writeTestKubeconfig(t)

	SetClusterContextName(contextA)
	require.True(t, IsConnectedToCluster(), "expected context-a to load successfully")
	cfgA := GetK8sConfig()
	require.NotNil(t, cfgA)
	require.Equal(t, "https://a.example.invalid:6443", cfgA.Host)
	require.Equal(t, contextA, GetContextName())

	SetClusterContextName(contextB)
	require.True(t, IsConnectedToCluster(), "expected context-b to load successfully")
	cfgB := GetK8sConfig()
	require.NotNil(t, cfgB)
	require.Equal(t, "https://b.example.invalid:6443", cfgB.Host,
		"GetK8sConfig still serving context-a's server after switching to context-b")
	require.Equal(t, contextB, GetContextName())
}

// Repeated calls with the SAME context must not pay a reload: this is the
// single-cluster-per-process case (the CLI's normal path), and it must keep
// whatever caching benefit the existing K8SConfig/clientConfigAPI cache provided.
func TestSetClusterContextName_SameContextIsNoop(t *testing.T) {
	resetClientState(t)
	t.Cleanup(func() { resetClientState(t) })

	contextA, _ := writeTestKubeconfig(t)

	SetClusterContextName(contextA)
	require.True(t, IsConnectedToCluster())
	cfgFirst := GetK8sConfig()
	require.NotNil(t, cfgFirst)

	SetClusterContextName(contextA)
	cfgSecond := GetK8sConfig()

	require.Same(t, cfgFirst, cfgSecond, "same context name should not invalidate the cached config")
}

// A load failure on the first context must not permanently suppress
// IsConnectedToCluster for a later, valid context: connectedToCluster is reset on
// every actual context change, not left latched to the previous context's outcome.
func TestSetClusterContextName_ResetsConnectedToClusterAcrossContexts(t *testing.T) {
	resetClientState(t)
	t.Cleanup(func() { resetClientState(t) })

	oldKubeconfig, hadKubeconfig := os.LookupEnv("KUBECONFIG")
	t.Cleanup(func() {
		if hadKubeconfig {
			_ = os.Setenv("KUBECONFIG", oldKubeconfig)
		} else {
			_ = os.Unsetenv("KUBECONFIG")
		}
	})

	// Point KUBECONFIG at a path that does not exist, so loading "context-a" fails
	// and connectedToCluster is driven to false.
	require.NoError(t, os.Setenv("KUBECONFIG", filepath.Join(t.TempDir(), fmt.Sprintf("missing-%d", os.Getpid()))))
	SetClusterContextName("context-a")
	require.False(t, IsConnectedToCluster(), "expected missing kubeconfig to fail to load")

	// Now point KUBECONFIG at a real, valid config and switch context.
	contextA, contextB := writeTestKubeconfig(t)
	_ = contextA
	SetClusterContextName(contextB)
	require.True(t, IsConnectedToCluster(),
		"context-a's load failure must not suppress a successful load for context-b")
}
