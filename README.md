# faas-k8s

OpenFaaS Kubernetes CRD &amp; Controller

### Local run

Create OpenFaaS CRD:
```bash
$ kubectl create -f artifacts/openfaas-crd.yaml
```

Start OpenFaaS controller (assumes you have a working GKE kubeconfig):
```bash
$ go run *.go -kubeconfig=$HOME/.kube/config -logtostderr=true
```

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
curl -d '{"service":"nodeinfo","image":"functions/nodeinfo:burner","envProcess":"node main.js","labels":{"com.openfaas.scale.min":"2","com.openfaas.scale.max":"15"},"environment":{"output":"verbose","debug":"true"}}' -X POST  http://localhost:9090/system/functions
```

List functions:

```bash
curl http://localhost:9090/system/functions | jq .
```

Scale PODs up/down:

```bash
curl -d '{"serviceName":"nodeinfo", "replicas": 3}' -X POST http://localhost:9090/system/scale-function/nodeinfo
```

Remove function:

```bash
curl -d '{"functionName":"nodeinfo"}' -X DELETE http://localhost:9090/system/functions
```

### Instrumentation

Prometheus route:

```bash
curl http://localhost:9090/metrics
```

Profiling is enabled by default, to disable it set `pprof` environment variable to `false`.

Pprof web UI can be access at `http://localhost:9090/debug/pprof/`. The goroutine, heap and threadcreate 
profilers are enabled along with the full goroutine stack dump.

Run the heap profiler:

```bash
go tool pprof goprofex http://localhost:9090/debug/pprof/heap
```

Run the goroutine profiler:

```bash
go tool pprof goprofex http://localhost:9090/debug/pprof/goroutine
```
