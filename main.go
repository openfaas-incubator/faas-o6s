package main

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	informers "github.com/openfaas-incubator/openfaas-operator/pkg/client/informers/externalversions"
	"github.com/openfaas-incubator/openfaas-operator/pkg/controller"
	"github.com/openfaas-incubator/openfaas-operator/pkg/server"
	"github.com/openfaas-incubator/openfaas-operator/pkg/signals"
	"github.com/openfaas-incubator/openfaas-operator/pkg/version"
	corev1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	masterURL  string
	kubeconfig string
)

var pullPolicyOptions = map[string]bool {
	"Always": true,
	"IfNotPresent": true,
	"Never": true,
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()

	sha, release := version.GetReleaseInfo()
	glog.Infof("Starting OpenFaaS controller version: %s commit: %s", release, sha)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building Kubernetes clientset: %s", err.Error())
	}

	faasClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building OpenFaaS clientset: %s", err.Error())
	}

	functionNamespace := "openfaas-fn"
	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	imagePullPolicy := corev1.PullAlways
	if val, exists := os.LookupEnv("image_pull_policy"); exists {
		if !pullPolicyOptions[val] {
			glog.Fatalf("Invalid image_pull_policy configured: %s", val)
		}
		imagePullPolicy = corev1.PullPolicy(val)
	}

	defaultResync := time.Second * 30

	kubeInformerOpt := kubeinformers.WithNamespace(functionNamespace)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, defaultResync, kubeInformerOpt)

	faasInformerOpt := informers.WithNamespace(functionNamespace)
	faasInformerFactory := informers.NewSharedInformerFactoryWithOptions(faasClient, defaultResync, faasInformerOpt)

	ctrl := controller.NewController(kubeClient, faasClient, kubeInformerFactory, faasInformerFactory, imagePullPolicy)

	go kubeInformerFactory.Start(stopCh)
	go faasInformerFactory.Start(stopCh)
	go server.Start(faasClient, kubeClient)

	if err = ctrl.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
