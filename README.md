# Creating Custom Operators

The `memhog-operator` is an example on how to create custom operators for
Kubernetes.

The purpose of the `memhog-operator` is to watch for Pods in a namespace and
monitor its memory usage. If the memory consumption of the Pod crosses a
threshold, it will be vertically autoscaled by the operator.

Specifically, the operator will deploy a new copy of the Pod with a higher
set of resource requests and limit, and then terminate the original Pod.
The details of the higher resources are held within an `AppMonitor`,
a custom TPR.

[memhog](https://github.com/metral/memhog): An example Pod that this operator would monitor.

> Note: The `memhog-operator` is strictly for demo purposes. It is not intended
to be used for any other use-cases.

## Operator Structure

The `memhog-operator` is a combination of a custom ThirdPartyResource (TPR)
known as the `AppMonitor`, and a custom controller to enforce state.

The `AppMonitor` encapsulates the autoscaling details for a Pod.
The controller watches a Namespace for an `AppMonitor`, and for Pods that wish 
to be monitored (via Annotation). It then applies the operational
thresholds and requirements declared in the `AppMonitor` onto the Pod.

## Process

* To monitor the Pod's resource memory consumption, the operator requires that the Pod have an annotation in its `spec.template.metadata` to associate itself with the `memhog-operator`.

  e.g. The [memhog](https://github.com/metral/memhog) example is annotated as such:

  <pre><code>apiVersion: extensions/v1beta1
  kind: Deployment
  metadata:
    name: memhog
  spec:
    replicas: 1
    template:
      metadata:
        labels:
          name: memhog
        <b>annotations:
          app-monitor.kubedemo.com/monitor: "true"</b>
      spec:
        containers:
        - name: memhog
          image: quay.io/metral/memhog:v0.0.1
          imagePullPolicy: Always
          resources:
            limits:
              memory: 384Mi
            requests:
              memory: 256Mi
          ...
  </code></pre>

* The Pod runs with a default set of resource requests & limits.

* In the Pod's namespace there must also be an instantiated object of the custom
`AppMonitor` TPR that the operator depends on e.g.:

  ```
  apiVersion: kubedemo.com/v1
  kind: AppMonitor
  metadata:
    name: mymonitor
  spec:
    memThresholdPercent: 75   # Percentage of (memory used) / (memory limit)
    memMultiplier: 2          # Multiplier factor used to increase memory resource requests & limits
  ```
* The `memhog-operator` will watch the Pod's memory usage by querying
Prometheus, apply the `AppMonitor` to the Pod as memory usage is retrieved, and redeploy the Pod 
with higher resource requests & limits if the `AppMonitor` thresholds are crossed.
* If the Pod is redeployed, it will have updated resource requests & limits
e.g.:
```
  ...
  resources:
    limits:
      memory: 768Mi
    requests:
      memory: 512Mi
  ...
```

### Building & Running

```
// Build
$ glide up -v
$ make
```

> Note: Prometheus is assumed to be running on http://localhost:9090
```
// Create cluster role & cluster role binding to work with TPR's.
$ kubectl create -f k8s/roles/role.yaml

// Run the operator
$ $GOPATH/bin/memhog-operator -v2 --prometheus-addr=http://localhost:9090 --kubeconfig=$HOME/.kube/config
```
