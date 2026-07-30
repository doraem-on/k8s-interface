package v1

import (
	"encoding/json"
	"strings"

	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/kubescape/k8s-interface/cloudsupport/mockobjects"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"google.golang.org/api/cloudresourcemanager/v1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func NewGKESupportMock() *GKESupportMock {
	return &GKESupportMock{}
}

type GKESupportMock struct {
}

// Get descriptive info about cluster running in GKE.
func (gkeSupportM *GKESupportMock) GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error) {
	c := &containerpb.Cluster{}
	err := json.Unmarshal([]byte(mockobjects.GkeDescriptor), c)
	return c, err
}

func (gkeSupportM *GKESupportMock) GetName(clusterDescribe *containerpb.Cluster) string {
	return clusterDescribe.Name
}

func (gkeSupportM *GKESupportMock) GetProject(cluster string) (string, error) {
	return "", nil
}

func (gkeSupportM *GKESupportMock) GetRegion(cluster string) (string, error) {
	return "", nil
}

func (gkeSupportM *GKESupportMock) GetContextName(cluster string) string {
	parsedName := strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return ""
	}
	clusterName := parsedName[3]
	if clusterName != "" {
		return clusterName
	}
	cluster = k8sinterface.GetContextName()
	parsedName = strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return ""
	}
	return parsedName[3]
}

func (gkeSupportM *GKESupportMock) GetIAMMappings(project string) (map[string]string, map[string]string, error) {
	return nil, nil, nil
}

func (gkeSupportM *GKESupportMock) GetDescribeRepositories(project string, region string) ([]*artifactregistrypb.Repository, error) {
	return []*artifactregistrypb.Repository{{Name: "mock-repo"}}, nil
}

func (gkeSupportM *GKESupportMock) GetListEntitiesForPolicies(project string) (*cloudresourcemanager.Policy, error) {
	return &cloudresourcemanager.Policy{
		Version: 3,
		Bindings: []*cloudresourcemanager.Binding{
			{
				Role: "roles/viewer",
			},
			{
				Role: "roles/editor",
				Condition: &cloudresourcemanager.Expr{
					Expression: "request.time < timestamp('2025-01-01T00:00:00Z')",
					Title:      "expires_end_of_2024",
				},
			},
		},
	}, nil
}
