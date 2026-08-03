package workloadinterface

import (
	"sort"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/workloadmethods/onesecretpercontainer.json
	oneSecretPerContainer string
	//go:embed testdata/workloadmethods/multiplesecretspercontainer.json
	multipleSecretsPerContainer string
	//go:embed testdata/workloadmethods/secretandconfigmapforcontainer.json
	secretAndConfigMapForContainer string
	//go:embed testdata/workloadmethods/secretandconfigmapforcontainercronjob.json
	secretAndConfigMapForContainerCronjob string
	//go:embed testdata/workloadmethods/secretandconfigmapforcontainerdeployment.json
	secretAndConfigMapForContainerDeployment string
	//go:embed testdata/workloadmethods/oneconfigmappercontainer.json
	oneConfigMapPerContainer string
	//go:embed testdata/workloadmethods/multipleconfigmapspercontainer.json
	multipleConfigMapsPerContainer string

	//go:embed testdata/workloadmethods/cfgMapsfromvolumesandenvfrom.json
	cfgMapsFromVolumesAndEnvFrom string

	//go:embed testdata/workloadmethods/cfgmapwithvaluefrom.json
	cfgMapWithValueFrom string

	//go:embed testdata/workloadmethods/secretsfromvolumesandenvfrom.json
	secretsFromVolumesAndEnvFrom string

	//go:embed testdata/workloadmethods/secretwithvaluefrom.json
	secretWithValueFrom string

	//go:embed testdata/workloadmethods/multiplesecretssamename.json
	multipleSecretsSameName string

	//go:embed testdata/workloadmethods/multipleconfigmapssamename.json
	multipleConfigMapsSameName string

	//go:embed testdata/workloadmethods/podstatusrunning.json
	podStatusRunning string

	//go:embed testdata/workloadmethods/podstatusnotpresent.json
	podStatusNotPresent string

	//go:embed testdata/workloadmethods/deploymentWithemptyLabels.json
	deploymentWithemptyLabels string

	//go:embed testdata/workloadmethods/podmountwithvolume.json
	podMountWithVolume string

	//go:embed testdata/workloadmethods/podmountnohostvolume.json
	podMountNoHostVolume string

	mockDeployment = `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	mockService    = `{"apiVersion":"v1","kind":"Service","metadata":{"creationTimestamp":"2021-12-06T14:01:16Z","labels":{"app":"armo-vuln-scan","app.kubernetes.io\/managed-by":"Helm"},"name":"armo-vuln-scan","resourceVersion":"351796","uid":"12bd4f9f-3ec6-4113-8ec6-0b8a1c772deb"},"spec":{"clusterIP":"10.107.7.78","clusterIPs":["10.107.7.78"],"internalTrafficPolicy":"Cluster","ipFamilies":["IPv4"],"ipFamilyPolicy":"SingleStack","ports":[{"port":8080,"protocol":"TCP","targetPort":8080}],"selector":{"app":"armo-vuln-scan"},"sessionAffinity":"None","type":"ClusterIP"},"status":{"loadBalancer":{}}}`
)

func TestLabels(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Error(err)
	}
	if workload.GetKind() != "Deployment" {
		t.Errorf("wrong kind")
	}
	if workload.GetNamespace() != "default" {
		t.Errorf("wrong namespace")
	}
	if workload.GetName() != "demoservice-server" {
		t.Errorf("wrong name")
	}
}

func TestGetReplicas(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Error(err)
	}
	SetInMap(workload.GetObject(), []string{"spec"}, "replicas", 3)
	assert.Equal(t, 3, workload.GetReplicas())
}
func TestGetSelector(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Error(err)
	}
	l, err := workload.GetSelector()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(l.MatchLabels))
	assert.Equal(t, 0, len(l.MatchExpressions))
}

func TestGetServiceSelector(t *testing.T) {
	workload, err := NewWorkload([]byte(mockService))
	if err != nil {
		t.Error(err)
	}

	s := workload.GetServiceSelector()
	expected := map[string]string{
		"app": "armo-vuln-scan",
	}

	assert.Equal(t, expected, s)

}

func TestSetNamespace(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"demoservice-server"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"demoservice-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	workload.SetNamespace("default")
	if workload.GetNamespace() != "default" {
		t.Errorf("wrong namespace")
	}
}
func TestSetLabels(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	workload.SetLabel("bla", "daa")
	v, ok := workload.GetLabel("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find label")
	}
	workload.RemoveLabel("bla")
	v2, ok2 := workload.GetLabel("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("label not deleted")
	}
}

func TestSetAnnotations(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	workload.SetAnnotation("bla", "daa")
	v, ok := workload.GetAnnotation("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find annotation")
	}
	workload.RemoveAnnotation("bla")
	v2, ok2 := workload.GetAnnotation("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("annotation not deleted")
	}
}
func TestSetPodLabels(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	workload.SetPodLabel("bla", "daa")
	v, ok := workload.GetPodLabel("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find label")
	}
	workload.RemovePodLabel("bla")
	v2, ok2 := workload.GetPodLabel("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("label not deleted")
	}
}

func TestGetResourceVersion(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	if workload.GetResourceVersion() != "1016043" {
		t.Errorf("wrong resourceVersion")
	}

}
func TestGetUID(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Error(err)
	}
	if workload.GetUID() != "e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e" {
		t.Errorf("wrong UID")
	}

}
func TestGetID(t *testing.T) {
	dWorkload, _ := NewWorkload([]byte(mockDeployment))
	assert.Equal(t, "apps/v1/default/Deployment/demoservice-server", dWorkload.GetID())

	sWorkload, _ := NewWorkload([]byte(mockService))
	assert.Equal(t, "/v1//Service/armo-vuln-scan", sWorkload.GetID())

}

func TestGetSecretsOfContainer(t *testing.T) {
	tests := []struct {
		name          string
		want          map[string][]string
		responseError error
		mockData      string
	}{
		{
			name: "One secret per container",
			want: map[string][]string{
				"container1": {"secret1"},
				"container2": {"secret2"},
			},
			responseError: nil,
			mockData:      oneSecretPerContainer,
		},
		{
			name: "Multiple secrets per container",
			want: map[string][]string{
				"container1": {"secret1", "secret2"},
				"container2": {"secret3", "secret4"},
			},
			responseError: nil,
			mockData:      multipleSecretsPerContainer,
		},
		{
			name: "Secret and configmap per container",
			want: map[string][]string{
				"container1": {"secret1", "secret2"},
				"container2": {"secret3", "secret4"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainer,
		},
		{
			name: "Secret and configmap per container - cronjob",
			want: map[string][]string{
				"container1": {"secret1", "secret2"},
				"container2": {"secret3", "secret4"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainerCronjob,
		},
		{
			name: "Secret and configmap per container - deployment",
			want: map[string][]string{
				"container1": {"secret1", "secret2"},
				"container2": {"secret3", "secret4"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainerDeployment,
		},
		{
			name: "Secret from Volumes and envFrom",
			want: map[string][]string{
				"container1": {"secret1", "secret2", "special-secret"},
			},
			responseError: nil,
			mockData:      secretsFromVolumesAndEnvFrom,
		},
		{
			name: "Secret with valueFrom",
			want: map[string][]string{
				"mycontainer": {"mysecret"},
			},
			responseError: nil,
			mockData:      secretWithValueFrom,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			secretsOfContainer, err := workload.GetSecretsOfContainer()
			assert.Equal(t, err, tc.responseError)
			assert.Equal(t, tc.want, secretsOfContainer)
		})
	}

}

func TestGetConfigMapsOfContainer(t *testing.T) {
	tests := []struct {
		name          string
		want          map[string][]string
		responseError error
		mockData      string
	}{
		{
			name: "One configmap per container",
			want: map[string][]string{
				"container1": {"config1"},
				"container2": {"config2"},
			},
			responseError: nil,
			mockData:      oneConfigMapPerContainer,
		},
		{
			name: "Multiple configmaps per container",
			want: map[string][]string{
				"container1": {"config1", "config2"},
				"container2": {"config3", "config4"},
			},
			responseError: nil,
			mockData:      multipleConfigMapsPerContainer,
		},
		{
			name: "Secret and configmap per container",
			want: map[string][]string{
				"container1": {"config1"},
				"container2": {"config2"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainer,
		},
		{
			name: "Secret and configmap per container - cronjob",
			want: map[string][]string{
				"container1": {"config1"},
				"container2": {"config2"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainerCronjob,
		},
		{
			name: "Secret and configmap per container - deployment",
			want: map[string][]string{
				"container1": {"config1"},
				"container2": {"config2"},
				"container3": {},
			},
			responseError: nil,
			mockData:      secretAndConfigMapForContainerDeployment,
		},
		{
			name: "ConfigMap from Volumes and envFrom",
			want: map[string][]string{
				"container1": {"config1", "config2", "special-config"},
			},
			responseError: nil,
			mockData:      cfgMapsFromVolumesAndEnvFrom,
		},
		{
			name: "ConfigMap with valueFrom",
			want: map[string][]string{
				"test-container": {"special-config"},
			},
			responseError: nil,
			mockData:      cfgMapWithValueFrom,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			configMapsOfContainer, err := workload.GetConfigMapsOfContainer()
			assert.Equal(t, err, tc.responseError)
			assert.Equal(t, tc.want, configMapsOfContainer)
		})
	}

}

func TestGetSecrets(t *testing.T) {
	tests := []struct {
		name          string
		want          []string
		responseError error
		mockData      string
	}{
		{
			name:          "No secrets",
			want:          []string{},
			responseError: nil,
			mockData:      oneConfigMapPerContainer,
		},
		{
			name:          "Multiple secrets",
			want:          []string{"secret1", "secret2", "secret3", "secret4"},
			responseError: nil,
			mockData:      secretAndConfigMapForContainer,
		},
		{
			name:          "Secret from Volumes and envFrom",
			want:          []string{"secret1", "secret2", "special-secret"},
			responseError: nil,
			mockData:      secretsFromVolumesAndEnvFrom,
		},
		{
			name:          "Secret with valueFrom",
			want:          []string{"mysecret"},
			responseError: nil,
			mockData:      secretWithValueFrom,
		},
		{
			name:          "Multiple secrets with same name",
			want:          []string{"secret1"},
			responseError: nil,
			mockData:      multipleSecretsSameName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			secrets, err := workload.GetSecrets()
			assert.Equal(t, err, tc.responseError)
			sort.Strings(tc.want)
			sort.Strings(secrets)
			assert.Equal(t, tc.want, secrets)
		})
	}

}

func TestGetConfigMaps(t *testing.T) {
	tests := []struct {
		name          string
		want          []string
		responseError error
		mockData      string
	}{
		{
			name:          "No configmaps",
			want:          []string{},
			responseError: nil,
			mockData:      oneSecretPerContainer,
		},
		{
			name:          "Multiple configmaps",
			want:          []string{"config1", "config2"},
			responseError: nil,
			mockData:      secretAndConfigMapForContainer,
		},
		{
			name:          "ConfigMap from Volumes and envFrom",
			want:          []string{"config1", "config2", "special-config"},
			responseError: nil,
			mockData:      cfgMapsFromVolumesAndEnvFrom,
		},
		{
			name:          "ConfigMap with valueFrom",
			want:          []string{"special-config"},
			responseError: nil,
			mockData:      cfgMapWithValueFrom,
		},
		{
			name:          "Multiple configmaps with same name",
			want:          []string{"config1"},
			responseError: nil,
			mockData:      multipleConfigMapsSameName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			configMaps, err := workload.GetConfigMaps()
			assert.Equal(t, err, tc.responseError)
			sort.Strings(tc.want)
			sort.Strings(configMaps)
			assert.Equal(t, tc.want, configMaps)
		})
	}

}

func TestGetPodStatus(t *testing.T) {
	tests := []struct {
		name          string
		want          string
		responseError error
		mockData      string
	}{
		{
			name:          "Pod status running",
			want:          "Running",
			responseError: nil,
			mockData:      podStatusRunning,
		},
		{
			name:          "Pod status not present",
			want:          "",
			responseError: nil,
			mockData:      podStatusNotPresent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			podStatus, err := workload.GetPodStatus()
			assert.Equal(t, string(podStatus.Phase), tc.want)
			assert.Equal(t, err, tc.responseError)
		})
	}

}

func TestGetLabels(t *testing.T) {
	tests := []struct {
		name     string
		want     map[string]string
		mockData string
	}{
		{
			name:     "null labels will not be returned",
			want:     map[string]string{"app": "my-app", "version": "1.0"},
			mockData: deploymentWithemptyLabels,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}
			labels := workload.GetLabels()
			assert.Equal(t, tc.want, labels)
		})
	}

}

func TestGetSpecPath(t *testing.T) {

	tests := []struct {
		name     string
		mockData string
		want     string
	}{
		{
			name:     "deployment",
			mockData: secretAndConfigMapForContainerDeployment,
			want:     "spec.template.spec",
		},
		{
			name:     "cronjob",
			mockData: secretAndConfigMapForContainerCronjob,
			want:     "spec.jobTemplate.spec.template.spec",
		},
		{
			name:     "pod",
			mockData: podStatusRunning,
			want:     "spec",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			workload, err := NewWorkload([]byte(tc.mockData))
			if err != nil {
				t.Error(err)
			}

			path, err := workload.GetSpecPath()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tc.want, path)
		})
	}

}

func TestGetVolumes(t *testing.T) {
	workload, err := NewWorkload([]byte(podMountWithVolume))
	if err != nil {
		t.Error(err)
	}

	volumes, err := workload.GetVolumes()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(volumes))

}

func TestAccessorsWithNonStringMetadata(t *testing.T) {
	// A manifest can be valid enough to be recognized as a workload and still
	// carry a non-string value where a string is expected. The accessors used
	// to assert the type outright and panic on such a document.
	w, err := NewWorkload([]byte(`{
		"apiVersion": ["oops"],
		"kind": 1,
		"metadata": {
			"namespace": ["oops"],
			"name": 42,
			"generateName": ["oops"],
			"resourceVersion": 7,
			"uid": ["oops"],
			"labels": {"app": ["nope"]},
			"annotations": {"note": 7}
		},
		"spec": {
			"selector": {"matchLabels": {"app": ["nope"], "tier": "web"}},
			"template": {
				"metadata": {
					"labels": {"app": ["nope"]},
					"annotations": {"note": 7}
				},
				"spec": {"serviceAccountName": ["oops"]}
			}
		}
	}`))
	assert.NoError(t, err)

	for _, tc := range []struct {
		name string
		got  func(*Workload) string
	}{
		{"GetNamespace", (*Workload).GetNamespace},
		{"GetName", (*Workload).GetName},
		{"GetApiVersion", (*Workload).GetApiVersion},
		{"GetKind", (*Workload).GetKind},
		{"GetGenerateName", (*Workload).GetGenerateName},
		{"GetResourceVersion", (*Workload).GetResourceVersion},
		{"GetUID", (*Workload).GetUID},
		{"GetServiceAccountName", (*Workload).GetServiceAccountName},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, "", tc.got(w))
		})
	}

	assert.NotPanics(t, func() { _ = w.GetID() })

	_, ok := w.GetLabel("app")
	assert.False(t, ok)
	_, ok = w.GetAnnotation("note")
	assert.False(t, ok)
	_, ok = w.GetPodLabel("app")
	assert.False(t, ok)
	_, ok = w.GetPodAnnotation("note")
	assert.False(t, ok)

	// the plural accessors read the same fields
	assert.Equal(t, map[string]string{}, w.GetLabels())
	assert.Equal(t, map[string]string{}, w.GetPodLabels())
	assert.Equal(t, map[string]string{"note": "7"}, w.GetAnnotations())
	assert.Equal(t, map[string]string{"note": "7"}, w.GetPodAnnotations())

	// a mixed matchLabels keeps the string entry and drops the rest
	selector, err := w.GetSelector()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"tier": "web"}, selector.MatchLabels)
}

// TestContainersWithNonMapValues covers the case where the container itself is
// not a map, rather than one of its values.
func TestContainersWithNonMapValues(t *testing.T) {
	w, err := NewWorkload([]byte(`{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "x", "labels": ["oops"], "annotations": "oops"},
		"data": "oops",
		"spec": {
			"selector": ["oops"],
			"template": {"metadata": {"labels": "oops", "annotations": ["oops"]}}
		}
	}`))
	assert.NoError(t, err)

	assert.Nil(t, w.GetLabels())
	assert.Nil(t, w.GetAnnotations())
	assert.Nil(t, w.GetPodLabels())
	assert.Nil(t, w.GetPodAnnotations())
	assert.Nil(t, w.GetData())
	assert.Nil(t, w.GetServiceSelector())
}

// TestSettersWithNonMapNode covers the write path: SetInMap used to test only
// for presence, so an existing non-map node left it writing into a nil map.
func TestSettersWithNonMapNode(t *testing.T) {
	w, err := NewWorkload([]byte(`{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": "oops"
	}`))
	assert.NoError(t, err)

	assert.NotPanics(t, func() { w.SetNamespace("default") })
	assert.Equal(t, "default", w.GetNamespace())

	assert.NotPanics(t, func() { w.SetLabel("app", "web") })
	label, ok := w.GetLabel("app")
	assert.True(t, ok)
	assert.Equal(t, "web", label)

	// an existing, valid map must not be replaced
	w2, err := NewWorkload([]byte(`{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {"name": "keep-me"}
	}`))
	assert.NoError(t, err)
	w2.SetNamespace("default")
	assert.Equal(t, "keep-me", w2.GetName())
	assert.Equal(t, "default", w2.GetNamespace())

	assert.NotPanics(t, func() { w.RemoveLabel("app") })
}

func TestIsTypeListWorkloadsWithNonStringKind(t *testing.T) {
	assert.False(t, IsTypeListWorkloads(map[string]interface{}{"kind": 123}))
}
