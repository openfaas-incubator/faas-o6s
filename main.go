package main

import (
	"context"
	"flag"
	"k8s.io/client-go/tools/record"
	"os"
	"time"

	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	informers "github.com/openfaas-incubator/openfaas-operator/pkg/client/informers/externalversions"
	"github.com/openfaas-incubator/openfaas-operator/pkg/controller"
	"github.com/openfaas-incubator/openfaas-operator/pkg/server"
	"github.com/openfaas-incubator/openfaas-operator/pkg/signals"
	"github.com/openfaas-incubator/openfaas-operator/pkg/version"
	"github.com/openfaas/faas-netes/k8s"
	"github.com/openfaas/faas-netes/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	glog "k8s.io/klog"

	// required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	masterURL  string
	kubeconfig string
)

var pullPolicyOptions = map[string]bool{
	"Always":       true,
	"IfNotPresent": true,
	"Never":        true,
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	setupLogging()

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

	readConfig := types.ReadConfig{}
	osEnv := types.OsEnv{}
	config := readConfig.Read(osEnv)

	deployConfig := k8s.DeploymentConfig{
		RuntimeHTTPPort: 8080,
		HTTPProbe:       config.HTTPProbe,
		SetNonRootUser:  config.SetNonRootUser,
		ReadinessProbe: &k8s.ProbeConfig{
			InitialDelaySeconds: int32(config.ReadinessProbeInitialDelaySeconds),
			TimeoutSeconds:      int32(config.ReadinessProbeTimeoutSeconds),
			PeriodSeconds:       int32(config.ReadinessProbePeriodSeconds),
		},
		LivenessProbe: &k8s.ProbeConfig{
			InitialDelaySeconds: int32(config.LivenessProbeInitialDelaySeconds),
			TimeoutSeconds:      int32(config.LivenessProbeTimeoutSeconds),
			PeriodSeconds:       int32(config.LivenessProbePeriodSeconds),
		},
		ImagePullPolicy: config.ImagePullPolicy,
	}

	factory := controller.NewFunctionFactory(kubeClient, deployConfig)

	functionNamespace := "openfaas-fn"
	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	if !pullPolicyOptions[config.ImagePullPolicy] {
		glog.Fatalf("Invalid image_pull_policy configured: %s", config.ImagePullPolicy)
	}

	defaultResync := time.Second * 30

	kubeInformerOpt := kubeinformers.WithNamespace(functionNamespace)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, defaultResync, kubeInformerOpt)

	faasInformerOpt := informers.WithNamespace(functionNamespace)
	faasInformerFactory := informers.NewSharedInformerFactoryWithOptions(faasClient, defaultResync, faasInformerOpt)

	ctrl := controller.NewController(
		kubeClient,
		faasClient,
		kubeInformerFactory,
		faasInformerFactory,
		factory,
	)

	go kubeInformerFactory.Start(stopCh)
	go faasInformerFactory.Start(stopCh)
	go server.Start(faasClient, kubeClient, kubeInformerFactory)

	// leader election context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// cancel leader election context on shutdown signals
	go func() {
		<-stopCh
		cancel()
	}()

	runController := func() {
		if err = ctrl.Run(1, stopCh); err != nil {
			glog.Fatalf("Error running controller: %s", err.Error())
		}
	}

	enableLeaderElection := false
	if val, exists := os.LookupEnv("enable_leader_election"); exists && val == "true" {
		enableLeaderElection = true
	}

	if enableLeaderElection {
		startLeaderElection(ctx, runController, functionNamespace, kubeClient)
	} else {
		runController()
	}
}

func setupLogging() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	glog.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})
}

func startLeaderElection(ctx context.Context, run func(), ns string, kubeClient kubernetes.Interface) {
	configMapName := "openfaas-leader-election"
	id, err := os.Hostname()
	if err != nil {
		glog.Fatalf("Error running controller: %v", err)
	}
	id = id + "_" + string(uuid.NewUUID())

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.V(4).Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "openfaas-operator"})

	lock, err := resourcelock.New(
		resourcelock.ConfigMapsResourceLock,
		ns,
		configMapName,
		kubeClient.CoreV1(),
		resourcelock.ResourceLockConfig{
			EventRecorder: recorder,
			Identity:      id,
		},
	)
	if err != nil {
		glog.Fatalf("Error running controller: %v", err)
	}

	glog.Infof("Starting leader election id: %s configmap: %s namespace: %s", id, configMapName, ns)
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 60 * time.Second,
		RenewDeadline: 15 * time.Second,
		RetryPeriod:   5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				glog.Info("Acting as elected leader")
				run()
			},
			OnStoppedLeading: func() {
				glog.Infof("Leadership lost")
				os.Exit(1)
			},
			OnNewLeader: func(identity string) {
				if identity != id {
					glog.Infof("Another instance has been elected as leader: %v", identity)
				}
			},
		},
	})
}
