package operator

import (
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
)

// #############################################################################

// Note: The following code is custom settings particular to the new TPR in the
// cluster.

// #############################################################################

const TPRName string = "app-monitor"
const ResourceName string = "appmonitor"
const ResourceNamePlural string = "appmonitors"
const Domain string = "kubedemo.com"
const Version string = "v1"
const Description string = "Allow user to create an app monitor to supervise the app and the resources it needs."
const AppMonitorAnnotation string = "app-monitor.kubedemo.com/monitor"
const AppMonitorAnnotationRedeployInProgress string = "app-monitor.kubedemo.com/redeploy-in-progress"

// An AppMonitor redeploys an app when resource limits are exceeded.
type AppMonitorSpec struct {
	MemThresholdPercent float64 `json:memThresholdPercent`
	MemMultiplier       float64 `json:memMultiplier`
}

// #############################################################################

// Note: The following code is boilerplate code needed to satisfy the
// AppMontior as a resource in the cluster in terms of how it expects TPR's to
// be created, operate and used.

// #############################################################################

type AppMonitor struct {
	// TODO: add Name field for AppMonitor as its currently missing
	unversioned.TypeMeta `json:",inline"`
	Metadata             api.ObjectMeta `json:"metadata"`

	Spec AppMonitorSpec `json:"spec"`
}

type AppMonitorList struct {
	unversioned.TypeMeta `json:",inline"`
	Metadata             unversioned.ListMeta `json:"metadata"`

	Items []AppMonitor `json:"items"`
}
