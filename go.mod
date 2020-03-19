module github.com/openfaas/openfaas-operator

go 1.13

require (
	github.com/google/go-cmp v0.3.0
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/openfaas/faas v0.0.0-20191125105239-365f459b3f3a
	github.com/openfaas/faas-netes v0.0.0-20200204113738-b12f1b6c368e
	github.com/openfaas/faas-provider v0.0.0-20200101101649-8f7c35975e1b
	github.com/prometheus/client_golang v0.9.2
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	k8s.io/code-generator v0.17.4
	k8s.io/klog v1.0.0
)

// Pin the Kubernetes version to prevent faas-netes downgrading the packages
replace (
	k8s.io/api => k8s.io/api v0.17.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4
	k8s.io/client-go => k8s.io/client-go v0.17.4
	k8s.io/code-generator => k8s.io/code-generator v0.17.4
)
