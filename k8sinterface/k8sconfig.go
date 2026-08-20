package k8sinterface

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	logger "github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// DO NOT REMOVE - load cloud providers auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// stateMu guards every package-level variable below it: they record which
// cluster/context this package is currently talking to, and are read and
// written from the exported functions in this file (and
// GetK8SServerGitVersion in k8sdiscovery.go). Without a lock, concurrent
// calls into this package's API - e.g. two goroutines each driving their own
// scan against a different context, as in kubescape's sequential
// multi-context ("fleet scan") support - are a data race: go test -race
// flags concurrent SetClusterContextName/GetContextName/IsConnectedToCluster
// calls immediately (see #163).
//
// The functions below that touch this state are structured so the mutex is
// only ever held by non-reentrant leaf operations: several exported
// functions call each other (GetConfig calls SetClientConfigAPI,
// IsConnectedToCluster calls LoadK8sConfig and SetConnectedToCluster, ...),
// and sync.Mutex/RWMutex are not reentrant, so a naive Lock()/defer Unlock()
// at the top of every exported function would deadlock the first time one
// of them called another. Internal getX/setX helpers do the locking; the
// exported functions call those instead of touching the variables or each
// other while holding the lock. SetClusterContextName is the one exception:
// its cache-invalidation logic is a single atomic check-then-update across
// four variables, so it takes the lock directly for its whole body rather
// than composing single-variable helpers.
//
// These variables stay exported for source compatibility, and to avoid
// that, prefer the accessor functions in this file over reading/writing
// them directly.
var stateMu sync.RWMutex

var connectedToCluster = true
var clusterContextName = ""
var ConfigClusterServerName = ""
var K8SGitServerVersion = ""

// K8SConfig pointer to k8s config
var K8SConfig *restclient.Config
var clientConfigAPI *clientcmdapi.Config

func getConnectedToCluster() bool {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return connectedToCluster
}

func getClusterContextName() string {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return clusterContextName
}

func getK8SConfig() *restclient.Config {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return K8SConfig
}

func setK8SConfig(c *restclient.Config) {
	stateMu.Lock()
	defer stateMu.Unlock()
	K8SConfig = c
}

func getClientConfigAPILocked() *clientcmdapi.Config {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return clientConfigAPI
}

func setRunningIncluster(v bool) {
	stateMu.Lock()
	defer stateMu.Unlock()
	RunningIncluster = v
}

func getK8SGitServerVersion() string {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return K8SGitServerVersion
}

func setK8SGitServerVersion(v string) {
	stateMu.Lock()
	defer stateMu.Unlock()
	K8SGitServerVersion = v
}

// KubernetesApi -
type KubernetesApi struct {
	ApiExtensionsClient clientset.Interface
	KubernetesClient    kubernetes.Interface
	DynamicClient       dynamic.Interface
	DiscoveryClient     discovery.DiscoveryInterface
	Context             context.Context
	K8SConfig           *restclient.Config
}

// NewKubernetesApi -
func NewKubernetesApi() *KubernetesApi {
	var kubernetesClient *kubernetes.Clientset
	var err error

	if !IsConnectedToCluster() {
		logger.L().Fatal("failed to load kubernetes config: no configuration has been provided, try setting KUBECONFIG environment variable")
	}

	k8sConfig := GetK8sConfig()

	kubernetesClient, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new kubernetes client", helpers.Error(err))
	}

	dynamicClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new dynamic client", helpers.Error(err))
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new discovery client", helpers.Error(err))
	}

	apiExtensionsClient, err := clientset.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new discovery client", helpers.Error(err))
	}

	restclient.SetDefaultWarningHandler(restclient.NoWarnings{})
	InitializeMapResources(discoveryClient)

	return &KubernetesApi{
		ApiExtensionsClient: apiExtensionsClient,
		KubernetesClient:    kubernetesClient,
		DynamicClient:       dynamicClient,
		DiscoveryClient:     discoveryClient,
		Context:             context.Background(),
		K8SConfig:           k8sConfig,
	}
}
func (k8sAPI *KubernetesApi) GetKubernetesClient() kubernetes.Interface {
	return k8sAPI.KubernetesClient
}
func (k8sAPI *KubernetesApi) GetDynamicClient() dynamic.Interface {
	return k8sAPI.DynamicClient
}
func (k8sAPI *KubernetesApi) GetDiscoveryClient() discovery.DiscoveryInterface {
	return k8sAPI.DiscoveryClient
}

// RunningIncluster whether running in cluster
var RunningIncluster bool

// LoadK8sConfig load config from local file or from cluster
func LoadK8sConfig() error {
	kubeconfig, err := config.GetConfigWithContext(getClusterContextName())
	if err != nil {
		return fmt.Errorf("failed to load kubernetes config: %s", strings.ReplaceAll(err.Error(), "KUBERNETES_MASTER", "KUBECONFIG"))
	}
	if _, err := restclient.InClusterConfig(); err == nil {
		setRunningIncluster(true)
	} else {
		setRunningIncluster(false)
	}

	setK8SConfig(kubeconfig)
	return nil
}

// GetK8sConfig get config. load if not loaded yet
func GetK8sConfig() *restclient.Config {
	if !IsConnectedToCluster() {
		return nil
	}
	return getK8SConfig()
}

func GetContext() *clientcmdapi.Context {
	kubeConfig := GetConfig()
	if kubeConfig == nil {
		return nil
	}

	contextName := getClusterContextName()
	if contextName == "" {
		// if context name is not set, use the current context
		contextName = kubeConfig.CurrentContext
	}

	if context, exist := kubeConfig.Contexts[contextName]; exist && context != nil {
		// return the context
		return context
	}
	return nil
}

func SetClientConfigAPI(conf *clientcmdapi.Config) {
	stateMu.Lock()
	defer stateMu.Unlock()
	clientConfigAPI = conf
}

func IsConnectedToCluster() bool {
	if getK8SConfig() == nil {
		if err := LoadK8sConfig(); err != nil {
			SetConnectedToCluster(false)
		}
	}
	return getConnectedToCluster()
}

func GetContextName() string {
	if name := getClusterContextName(); name != "" {
		return name
	}

	if config := GetConfig(); config != nil {
		return config.CurrentContext
	}

	return ""
}

// get config from ~/.kube/config
func GetConfig() *clientcmdapi.Config {

	if !getConnectedToCluster() {
		return nil
	}
	if conf := getClientConfigAPILocked(); conf != nil {
		return conf
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{CurrentContext: getClusterContextName()})
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil
	}

	// set the config to the global variable
	SetClientConfigAPI(&config)

	return &config
}

// GetDefaultNamespace returns the default namespace for the current context
func GetDefaultNamespace() string {

	if context := GetContext(); context != nil {
		return context.Namespace
	}

	// return default namespace in case the context is not available
	return "default"
}

// GetCluster returns a pointer to the clientcmdapi Cluster object of the current context
func GetCluster() *clientcmdapi.Cluster {
	config := GetConfig()
	if config == nil {
		return nil
	}

	if context, exist := config.Clusters[GetContextName()]; exist && context != nil {
		// return the cluster as based on the context
		return context
	}
	return nil

}

// SetClusterContextName set the name of desired cluster context. The package will use this name when loading the context.
//
// If contextName differs from the currently set context, the cached client config
// (K8SConfig, clientConfigAPI) is invalidated so the next GetK8sConfig/GetConfig call
// reloads it for the new context instead of silently continuing to serve the previous
// context's config. Without this, GetContextName() reports the new context while
// GetK8sConfig() keeps returning the old one, since both LoadK8sConfig and GetConfig
// only (re)load when their respective cache is nil (see #158). connectedToCluster is
// reset to true so a load failure on the previous context does not suppress a retry
// for the new one; IsConnectedToCluster re-derives it from the new load attempt.
//
// Repeated calls with the same contextName are a no-op here, same as before, so the
// common single-cluster-per-process case keeps its caching benefit unchanged.
func SetClusterContextName(contextName string) {
	stateMu.Lock()
	defer stateMu.Unlock()
	if contextName != clusterContextName {
		K8SConfig = nil
		clientConfigAPI = nil
		connectedToCluster = true
	}
	clusterContextName = contextName
}

func SetK8SGitServerVersion(K8SGitServerVersionInput string) {
	setK8SGitServerVersion(K8SGitServerVersionInput)
}

func SetConfigClusterServerName(contextName string) {
	stateMu.Lock()
	defer stateMu.Unlock()
	ConfigClusterServerName = contextName
}

// GetK8sConfigClusterServerName get the server name of desired cluster context
func GetK8sConfigClusterServerName() string {

	config := GetConfig()
	if config == nil {
		return ""
	}

	if context, exist := config.Clusters[GetContextName()]; exist && context != nil {
		// return the server name of the context
		return context.Server
	}

	// return current context in case the server name is not available
	stateMu.RLock()
	defer stateMu.RUnlock()
	return ConfigClusterServerName
}

func SetConnectedToCluster(isConnected bool) {
	stateMu.Lock()
	defer stateMu.Unlock()
	connectedToCluster = isConnected
}
