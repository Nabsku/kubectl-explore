package explore_test

import (
	"testing"

	"github.com/keisku/kubectl-explore/explore"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
)

func TestOptions_serverPreferredResources(t *testing.T) {
	tests := []struct {
		name             string
		apiVersion       string
		resources        []*v1.APIResourceList
		expectedCount    int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:       "no api version filter returns all resources",
			apiVersion: "",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
					},
				},
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:       "filter by v1 returns only v1 resources",
			apiVersion: "v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:       "filter by apps/v1 returns only apps/v1 resources",
			apiVersion: "apps/v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:       "filter by non-existent api version returns error",
			apiVersion: "nonexistent/v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
					},
				},
			},
			expectedCount:    0,
			expectedError:    true,
			expectedErrorMsg: "no resources found for API version \"nonexistent/v1\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake discovery client with test resources
			fakeDiscovery := cmdtesting.NewFakeCachedDiscoveryClient()
			fakeDiscovery.PreferredResources = tt.resources

			// Create Options instance
			opts := explore.NewOptions(genericclioptions.IOStreams{})
			
			// Use helper functions to set private fields for testing
			// In production, these would be set via Complete() method
			explore.SetDiscoveryClient(opts, fakeDiscovery)
			explore.SetAPIVersion(opts, tt.apiVersion)

			// Test serverPreferredResources method
			result, err := explore.CallServerPreferredResources(opts)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Len(t, result, tt.expectedCount)
				
				// Verify that the returned resources match the expected API version filter
				if tt.apiVersion != "" {
					for _, list := range result {
						require.Equal(t, tt.apiVersion, list.GroupVersion)
					}
				}
			}
		})
	}
}

func TestOptions_listGVRs_withAPIVersionFilter(t *testing.T) {
	tests := []struct {
		name             string
		apiVersion       string
		resources        []*v1.APIResourceList
		expectedGVRs     []schema.GroupVersionResource
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:       "filter by v1 returns only v1 GVRs",
			apiVersion: "v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
						{Name: "services", Kind: "Service"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{
				{Group: "", Version: "v1", Resource: "pods"},
				{Group: "", Version: "v1", Resource: "services"},
			},
			expectedError: false,
		},
		{
			name:       "filter by apps/v1 returns only apps/v1 GVRs",
			apiVersion: "apps/v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []v1.APIResource{
						{Name: "deployments", Kind: "Deployment"},
						{Name: "replicasets", Kind: "ReplicaSet"},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{
				{Group: "apps", Version: "v1", Resource: "deployments"},
				{Group: "apps", Version: "v1", Resource: "replicasets"},
			},
			expectedError: false,
		},
		{
			name:       "filter by non-existent API version returns error",
			apiVersion: "nonexistent/v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
			},
			expectedGVRs:     nil,
			expectedError:    true,
			expectedErrorMsg: "no resources found for API version \"nonexistent/v1\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake discovery client with test resources
			fakeDiscovery := cmdtesting.NewFakeCachedDiscoveryClient()
			fakeDiscovery.PreferredResources = tt.resources

			// Create Options instance
			opts := explore.NewOptions(genericclioptions.IOStreams{})
			
			// Set private fields for testing
			explore.SetDiscoveryClient(opts, fakeDiscovery)
			explore.SetAPIVersion(opts, tt.apiVersion)

			// Test listGVRs method
			result, err := explore.CallListGVRs(opts)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Len(t, result, len(tt.expectedGVRs))
				
				// Verify the returned GVRs match expected ones
				for i, expectedGVR := range tt.expectedGVRs {
					require.Equal(t, expectedGVR, result[i])
				}
			}
		})
	}
}

func TestOptions_discover_withAPIVersionFilter(t *testing.T) {
	tests := []struct {
		name             string
		apiVersion       string
		resources        []*v1.APIResourceList
		expectedGVRCount int
		expectedMapKeys  []string  // Expected keys in the resource map
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:       "filter by autoscaling/v2 returns only autoscaling/v2 resources",
			apiVersion: "autoscaling/v2",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
				{
					GroupVersion: "autoscaling/v2",
					APIResources: []v1.APIResource{
						{Name: "horizontalpodautoscalers", Kind: "HorizontalPodAutoscaler", SingularName: "horizontalpodautoscaler", ShortNames: []string{"hpa"}},
					},
				},
			},
			expectedGVRCount: 1,
			// The map will contain entries for: name, kind, singular name, and short names
			expectedMapKeys:  []string{"horizontalpodautoscalers", "HorizontalPodAutoscaler", "horizontalpodautoscaler", "hpa"},
			expectedError:    false,
		},
		{
			name:       "filter by non-existent API version returns error",
			apiVersion: "nonexistent/v1",
			resources: []*v1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []v1.APIResource{
						{Name: "pods", Kind: "Pod"},
					},
				},
			},
			expectedGVRCount:  0,
			expectedMapKeys:   nil,
			expectedError:     true,
			expectedErrorMsg:  "no resources found for API version \"nonexistent/v1\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake discovery client with test resources
			fakeDiscovery := cmdtesting.NewFakeCachedDiscoveryClient()
			fakeDiscovery.PreferredResources = tt.resources

			// Create Options instance
			opts := explore.NewOptions(genericclioptions.IOStreams{})
			
			// Set private fields for testing
			explore.SetDiscoveryClient(opts, fakeDiscovery)
			explore.SetAPIVersion(opts, tt.apiVersion)

			// Test discover method
			resourceMap, gvrs, err := explore.CallDiscover(opts)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
				require.Nil(t, resourceMap)
				require.Nil(t, gvrs)
			} else {
				require.NoError(t, err)
				require.Len(t, gvrs, tt.expectedGVRCount)
				
				// Verify that the expected map keys are present
				if tt.expectedMapKeys != nil {
					for _, key := range tt.expectedMapKeys {
						require.Contains(t, resourceMap, key, "Expected key %s to be in resource map", key)
					}
				}
				
				// Verify that the GVRs match the expected API version
				if tt.apiVersion != "" {
					for _, gvr := range gvrs {
						expectedGV, _ := schema.ParseGroupVersion(tt.apiVersion)
						require.Equal(t, expectedGV.Group, gvr.Group)
						require.Equal(t, expectedGV.Version, gvr.Version)
					}
				}
			}
		})
	}
}

func TestNewCmd_APIVersionFlag(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedAPI string
	}{
		{
			name:        "no api-version flag sets empty string",
			args:        []string{},
			expectedAPI: "",
		},
		{
			name:        "api-version flag sets the value",
			args:        []string{"--api-version=apps/v1"},
			expectedAPI: "apps/v1",
		},
		{
			name:        "api-version flag with v1",
			args:        []string{"--api-version=v1"},
			expectedAPI: "v1",
		},
		{
			name:        "api-version flag with complex version",
			args:        []string{"--api-version=autoscaling/v2"},
			expectedAPI: "autoscaling/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := explore.NewCmd()
			
			// Set the args and parse flags
			cmd.SetArgs(tt.args)
			err := cmd.ParseFlags(tt.args)
			require.NoError(t, err)
			
			// Check if the flag was parsed correctly by inspecting the flag value
			flag := cmd.Flags().Lookup("api-version")
			require.NotNil(t, flag)
			require.Equal(t, tt.expectedAPI, flag.Value.String())
		})
	}
}