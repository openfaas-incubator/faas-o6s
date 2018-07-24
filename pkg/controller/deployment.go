package controller

import (
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// newDeployment creates a new Deployment for a Function resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Function resource that 'owns' it.
func newDeployment(
	function *faasv1.Function,
	existingSecrets map[string]*corev1.Secret,
	imagePullPolicy corev1.PullPolicy) *appsv1beta2.Deployment {

	envVars := makeEnvVars(function)
	labels := makeLabels(function)
	nodeSelector := makeNodeSelector(function.Spec.Constraints)
	livenessProbe := makeLivenessProbe()

	resources, err := makeResources(function)
	if err != nil {
		glog.Warningf("Function %s resources parsing failed: %v",
			function.Spec.Name, err)
	}

	deploymentSpec := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      function.Spec.Name,
			Namespace: function.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(function, schema.GroupVersionKind{
					Group:   faasv1.SchemeGroupVersion.Group,
					Version: faasv1.SchemeGroupVersion.Version,
					Kind:    faasKind,
				}),
			},
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: function.Spec.Replicas,
			Strategy: appsv1beta2.DeploymentStrategy{
				Type: appsv1beta2.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1beta2.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        function.Spec.Name,
					"controller": function.Name,
				},
			},
			RevisionHistoryLimit: int32p(5),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: map[string]string{"prometheus.io.scrape": "false"},
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []corev1.Container{
						{
							Name:  function.Spec.Name,
							Image: function.Spec.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: int32(functionPort), Protocol: corev1.ProtocolTCP},
							},
							ImagePullPolicy: imagePullPolicy,
							Env:             envVars,
							Resources:       *resources,
							LivenessProbe:   livenessProbe,
							ReadinessProbe:  livenessProbe,
						},
					},
				},
			},
		},
	}

	if err := UpdateSecrets(function, deploymentSpec, existingSecrets); err != nil {
		glog.Warningf("Function %s secrets update failed: %v",
			function.Spec.Name, err)
	}

	return deploymentSpec
}

func makeEnvVars(function *faasv1.Function) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if len(function.Spec.Handler) > 0 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "fprocess",
			Value: function.Spec.Handler,
		})
	}

	if function.Spec.Environment != nil {
		for k, v := range *function.Spec.Environment {
			envVars = append(envVars, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}
	}

	return envVars
}

func makeLabels(function *faasv1.Function) map[string]string {
	labels := map[string]string{
		"faas_function": function.Spec.Name,
		"app":           function.Spec.Name,
		"controller":    function.Name,
	}
	if function.Spec.Labels != nil {
		for k, v := range *function.Spec.Labels {
			labels[k] = v
		}
	}

	return labels
}

func makeLivenessProbe() *corev1.Probe {
	path := filepath.Join("/tmp/", ".lock")
	probe := &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{"cat", path},
			},
		},
		InitialDelaySeconds: 3,
		TimeoutSeconds:      1,
		PeriodSeconds:       5,
		SuccessThreshold:    1,
		FailureThreshold:    2,
	}

	return probe
}

func makeNodeSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")

			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

// deploymentNeedsUpdate determines if the function spec is different from the deployment spec
func deploymentNeedsUpdate(function *faasv1.Function, deployment *appsv1beta2.Deployment) bool {
	needsUpdate := false

	if function.Spec.Replicas != nil && *function.Spec.Replicas != *deployment.Spec.Replicas {
		glog.V(2).Infof("Function %s replica count changed from %d to %d",
			function.Spec.Name, *deployment.Spec.Replicas, *function.Spec.Replicas)
		needsUpdate = true
	}

	currentImage := deployment.Spec.Template.Spec.Containers[0].Image
	if function.Spec.Image != deployment.Spec.Template.Spec.Containers[0].Image {
		glog.V(2).Infof("Function %s image changed from %d to %d",
			function.Spec.Name, currentImage, function.Spec.Image)
		needsUpdate = true
	}

	currentEnv := deployment.Spec.Template.Spec.Containers[0].Env
	funcEnv := makeEnvVars(function)
	if envVarsNotEqual(currentEnv, funcEnv) {
		glog.V(2).Infof("Function %s envVars have changed",
			function.Spec.Name)
		needsUpdate = true
	}

	currentLabels := deployment.Spec.Template.Labels
	funcLabels := makeLabels(function)
	if labelsNotEqual(currentLabels, funcLabels) {
		glog.V(2).Infof("Function %s labels have changed",
			function.Spec.Name)
		needsUpdate = true
	}

	if secretsNotEqual(function.Spec.Secrets, deployment.Spec.Template.Spec.Volumes) {
		glog.V(2).Infof("Function %s secrets have changed",
			function.Spec.Name)
		needsUpdate = true
	}

	return needsUpdate
}

func envVarsNotEqual(a, b []corev1.EnvVar) bool {
	if len(a) != len(b) {
		return true
	}
	mb := map[string]bool{}
	for _, x := range b {
		mb[x.Name+x.Value] = true
	}

	for _, x := range a {
		if _, ok := mb[x.Name+x.Value]; !ok {
			return true
		}
	}
	return false
}

func labelsNotEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return true
	}
	mb := map[string]bool{}
	for v, x := range b {
		mb[v+x] = true
	}

	for v, x := range a {
		if _, ok := mb[v+x]; !ok {
			return true
		}
	}
	return false
}

func strArrayNotEqual(a, b []string) bool {
	if len(a) != len(b) {
		return true
	}

	for i, v := range a {
		if v != b[i] {
			return true
		}
	}
	return false
}

func secretsNotEqual(secrets []string, volumes []corev1.Volume) bool {
	if len(secrets) < 1 {
		return false
	}

	if len(secrets) > 0 && len(volumes) < 1 {
		return true
	}

	if len(secrets) > 0 && volumes[0].Projected == nil {
		return true
	}

	sources := []string{}
	for _, s := range volumes[0].Projected.Sources {
		if s.Secret != nil {
			sources = append(sources, s.Secret.Name)
		}
	}

	return strArrayNotEqual(secrets, sources)
}

func int32p(i int32) *int32 {
	return &i
}
