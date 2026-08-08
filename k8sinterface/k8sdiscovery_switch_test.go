package k8sinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

// stubDiscovery serves a canned resource list. Only ServerPreferredResources is
// exercised by InitializeMapResources, so the rest of the interface is left
// embedded and unimplemented.
type stubDiscovery struct {
	discovery.DiscoveryInterface
	resources []*metav1.APIResourceList
}

func (s *stubDiscovery) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return s.resources, nil
}

func clusterDiscovery(crdName string) *stubDiscovery {
	return &stubDiscovery{resources: []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "pods", Namespaced: true, Verbs: metav1.Verbs{"list"}},
			},
		},
		{
			GroupVersion: "example.io/v1",
			APIResources: []metav1.APIResource{
				{Name: crdName, Namespaced: true, Verbs: metav1.Verbs{"list"}},
			},
		},
	}}
}

// resetDiscoveryState clears the process-global discovery snapshot so each test
// starts from the same place.
func resetDiscoveryState(t *testing.T) {
	t.Helper()
	resourcesInfoLock.Lock()
	defer resourcesInfoLock.Unlock()
	resourceGroupMapping = map[string][]string{}
	resourceNamesapcedScope = []string{}
	ResourceClusterScope = []string{}
}

func hasResource(name string) bool {
	resourcesInfoLock.RLock()
	defer resourcesInfoLock.RUnlock()
	_, ok := resourceGroupMapping[name]
	return ok
}

func namespacedScopeHas(triplet string) bool {
	resourcesInfoLock.RLock()
	defer resourcesInfoLock.RUnlock()
	for _, r := range resourceNamesapcedScope {
		if r == triplet {
			return true
		}
	}
	return false
}

// A second live client, built for a different cluster, must not keep serving the
// first cluster's API discovery. InitializeMapResources returns early whenever
// the global snapshot is already populated, so today cluster B inherits A's.
func TestInitializeMapResources_ReplacesDiscoveryForSecondLiveClient(t *testing.T) {
	resetDiscoveryState(t)
	t.Cleanup(func() { resetDiscoveryState(t) })

	InitializeMapResources(clusterDiscovery("foos"))

	require.True(t, hasResource("foos"), "cluster A discovery should be loaded")
	require.False(t, hasResource("bars"), "cluster B CRD must not be present yet")

	// Switch to a client for a different cluster.
	InitializeMapResources(clusterDiscovery("bars"))

	assert.True(t, hasResource("bars"), "cluster B CRD must be discoverable after switching")
	assert.False(t, hasResource("foos"), "cluster A CRD must not survive the switch")
}

// The namespaced-scope slice has the same problem, and because setMapResources
// appends rather than replaces, a naive fix would leave a union of both clusters.
func TestInitializeMapResources_ReplacesNamespacedScopeForSecondLiveClient(t *testing.T) {
	resetDiscoveryState(t)
	t.Cleanup(func() { resetDiscoveryState(t) })

	InitializeMapResources(clusterDiscovery("foos"))
	require.True(t, namespacedScopeHas("example.io/v1/foos"))

	InitializeMapResources(clusterDiscovery("bars"))

	assert.True(t, namespacedScopeHas("example.io/v1/bars"), "cluster B resource must be namespaced-scoped")
	assert.False(t, namespacedScopeHas("example.io/v1/foos"), "cluster A resource must not survive the switch")
}

// Passing nil must keep the existing lazy/mock fallback behaviour.
func TestInitializeMapResources_NilClientStillFallsBackToMock(t *testing.T) {
	resetDiscoveryState(t)
	t.Cleanup(func() { resetDiscoveryState(t) })

	InitializeMapResources(nil)

	assert.True(t, hasResource("pods"), "the mock fallback should still populate discovery")
}
