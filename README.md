# openfaas-operator

[![Go Report Card](https://goreportcard.com/badge/github.com/openfaas-incubator/openfaas-operator)](https://goreportcard.com/report/github.com/openfaas-incubator/openfaas-operator) [![Build
Status](https://travis-ci.org/openfaas-incubator/openfaas-operator.svg?branch=master)](https://travis-ci.org/openfaas-incubator/openfaas-operator) [![GoDoc](https://godoc.org/github.com/openfaas-incubator/openfaas-operator?status.svg)](https://godoc.org/github.com/openfaas-incubator/openfaas-operator) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)

OpenFaaS Operator for Kubernetes 1.9 or newer

### Deploy

Deploy OpenFaaS Operator with Helm:

```bash
# create OpenFaaS namespaces
kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml

# add OpenFaaS Helm repo
helm repo add openfaas https://openfaas.github.io/faas-netes/

# get latest chart version and install
helm repo update && helm upgrade openfaas --install openfaas/openfaas \
    --namespace openfaas  \
    --set functionNamespace=openfaas-fn \
    --set operator.create=true
```

If you are upgrading from faas-netes you need to remove all functions and redeploy them after installing the operator.

Deploy a function with kubectl:

```bash
kubectl -n openfaas-fn apply -f ./artifacts/nodeinfo.yaml
```

On armhf use:

```bash
kubectl -n openfaas-fn apply -f ./artifacts/figlet-armhf.yaml
```

List functions, services, deployments and pods:

```bash
kubectl -n openfaas-fn get functions
kubectl -n openfaas-fn get all
``` 

Deploy a function with secrets:

```bash
kubectl -n openfaas-fn create secret generic faas-token --from-literal=faas-token=token
kubectl -n openfaas-fn create secret generic faas-key --from-literal=faas-key=key
kubectl -n openfaas-fn apply -f ./artifacts/gofast.yaml
```

Test that secrets are available inside the `gofast` pod:

```bash
kubectl -n openfaas-fn exec -it gofast-84fd464784-sd5ml -- sh

~ $ cat /var/openfaas/faas-key 
key

~ $ cat /var/openfaas/faas-token 
token
``` 

Test that node selectors work on GKE by adding the following to `gofast.yaml`:

```yaml
  constraints:
    - "cloud.google.com/gke-nodepool=default-pool"
```

Apply the function and check the deployment specs with:

```bash
kubectl -n openfaas-fn describe deployment gofast
```

#### Development build
The OpenFaaS Operator runs as a sidecar in the gateway pod. For end to end testing you need to update the sidecar to use
your development image.

1. Build,tag and push your image to your own public docker repository: 
i.e. `make build && docker tag openfaas/openfaas-operator:latest {user}/openfaas-operator:latest-dev`

2. Update your helm deployment
```
helm upgrade openfaas --install openfaas/openfaas \
    --namespace openfaas  \
    --set functionNamespace=openfaas-fn \
    --set operator.create=true \
--set operator.image={user}/openfaas-operator:latest-dev
```

### Local run

Create OpenFaaS CRD:
```bash
$ kubectl apply -f artifacts/operator-crd.yaml
```

Start OpenFaaS controller (assumes you have a working kubeconfig on the machine):

```bash
$ go build \
  && ./openfaas-operator -kubeconfig=$HOME/.kube/config -logtostderr=true -v=4
```

With `go run`

```bash
$ go run *.go -kubeconfig=$HOME/.kube/config -logtostderr=true
```

To use an alternative port set the `port` environmental variable to another value.

Create a function:
```bash
$ kubectl apply -f artifacts/nodeinfo.yaml
```

Check if nodeinfo deployment and service were created through the CRD:
```bash
$ kubectl get deployment nodeinfo
$ kubectl get service nodeinfo
```

Test if nodeinfo service can access the pods:
```bash
$ kubectl run -it --rm --restart=Never curl --image=byrnedo/alpine-curl --command -- sh
/ # curl -d 'verbose' http://nodeinfo.default:8080
```

Delete nodeinfo function:
```bash
kubectl delete -f artifacts/nodeinfo.yaml 
```

Check if nodeinfo pods, rc, deployment and service were removed:
```bash
kubectl get all
```

### API calls

Create or update a function:

```bash
curl -d '{"service":"nodeinfo","image":"functions/nodeinfo:burner","envProcess":"node main.js","labels":{"com.openfaas.scale.min":"2","com.openfaas.scale.max":"15"},"environment":{"output":"verbose","debug":"true"}}' -X POST  http://localhost:8081/system/functions
```

List functions:

```bash
curl -s http://localhost:8081/system/functions | jq .
```

Scale PODs up/down:

```bash
curl -d '{"serviceName":"nodeinfo", "replicas": 3}' -X POST http://localhost:8081/system/scale-function/nodeinfo
```

Get available replicas:

```bash
curl -s http://localhost:8081/system/function/nodeinfo | jq .availableReplicas
```

Remove function:

```bash
curl -d '{"functionName":"nodeinfo"}' -X DELETE http://localhost:8081/system/functions
```

### Logging

Verbosity levels:

* `-v=0` CRUD actions via API and Controller including errors
* `-v=2` function call duration (Proxy API)
* `-v=4` Kubernetes informers events (highly verbose)

### Instrumentation

Prometheus route:

```bash
curl http://localhost:8081/metrics
```

Profiling is disabled by default, to enable it set `pprof` environment variable to `true`.

Pprof web UI can be access at `http://localhost:8081/debug/pprof/`. The goroutine, heap and threadcreate 
profilers are enabled along with the full goroutine stack dump.

Run the heap profiler:

```bash
go tool pprof goprofex http://localhost:8081/debug/pprof/heap
```

Run the goroutine profiler:

```bash
go tool pprof goprofex http://localhost:8081/debug/pprof/goroutine
```
