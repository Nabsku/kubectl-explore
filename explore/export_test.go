package explore

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func SetDisablePrintPath(o *Options, b bool) {
	o.disablePrintPath = b
}

func SetShowBrackets(o *Options, b bool) {
	o.showBrackets = b
}

// Helper functions for testing GVR filtering functionality

func SetAPIVersion(o *Options, apiVersion string) {
	o.apiVersion = apiVersion
}

func SetDiscoveryClient(o *Options, discovery discovery.CachedDiscoveryInterface) {
	o.discovery = discovery
}

func CallServerPreferredResources(o *Options) ([]*v1.APIResourceList, error) {
	return o.serverPreferredResources()
}

func CallListGVRs(o *Options) ([]schema.GroupVersionResource, error) {
	return o.listGVRs()
}

func CallDiscover(o *Options) (map[string]*groupVersionAPIResource, []schema.GroupVersionResource, error) {
	return o.discover()
}
