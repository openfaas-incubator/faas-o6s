# faas-o6s

OpenFaaS Kubernetes CRD Controller

### Deploy

Deploy OpenFaaS with faas-netes:

```bash
git clone https://github.com/openfaas/faas-netes
cd faas-netes
kubectl apply -f ./namespaces.yml,./yaml
```

Deploy the CRD and faas-o6s controller in `openfaas` namespace:

```bash
# CRD
kubectl apply -f artifacts/o6s-crd.yaml
# RBAC
kubectl apply -f artifacts/o6s-rbac.yaml
# Service
kubectl apply -f artifacts/o6s-svc.yaml
# Deployment
kubectl apply -f artifacts/o6s-dep.yaml
```

Modify the gateway deployment and switch from faas-netes to faas-o6s:

```bash
      containers:
      - name: gateway
        image: functions/gateway:0.6.16
        imagePullPolicy: Always
        env:
        - name: functions_provider_url
          value: "http://faas-o6s.openfaas:8081/"
```

Deploy a function with kubectl:

```bash
kubectl -n openfaas-fn apply -f ./artifacts/nodeinfo.yaml
```

List functions, services, deployments and pods:

```bash
kubectl -n openfaas-fn get functions
kubectl -n openfaas-fn get all
``` 

### Local run

Create OpenFaaS CRD:
```bash
$ kubectl apply -f artifacts/o6s-crd.yaml
```

Start OpenFaaS controller (assumes you have a working kubeconfig on the machine):

```bash
$ go build \
  && ./faas-o6s -kubeconfig=$HOME/.kube/config -logtostderr=true
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
curl http://localhost:8081/system/functions | jq .
```

Scale PODs up/down:

```bash
curl -d '{"serviceName":"nodeinfo", "replicas": 3}' -X POST http://localhost:8081/system/scale-function/nodeinfo
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
