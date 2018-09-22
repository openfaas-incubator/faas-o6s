package controller

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-cmp/cmp"
	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	annotationFunctionSpec = "com.openfaas.function.spec"
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

	annotations := makeAnnotations(function)

	deploymentSpec := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        function.Spec.Name,
			Annotations: annotations,
			Namespace:   function.Namespace,
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
					Annotations: annotations,
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

	configureReadOnlyRootFilesystem(function, deploymentSpec)

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

func makeAnnotations(function *faasv1.Function) map[string]string {
	annotations := make(map[string]string)

	// disable scraping since the watchdog doesn't expose a metrics endpoint
	annotations["prometheus.io.scrape"] = "false"

	// copy function annotations
	if function.Spec.Annotations != nil {
		for k, v := range *function.Spec.Annotations {
			annotations[k] = v
		}
	}

	// save function spec in deployment annotations
	// used to detect changes in function spec
	specJSON, err := json.Marshal(function.Spec)
	if err != nil {
		glog.Errorf("Failed to marshal function spec: %s", err.Error())
		return annotations
	}

	annotations[annotationFunctionSpec] = string(specJSON)
	return annotations
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
	prevFnSpecJson := deployment.ObjectMeta.Annotations[annotationFunctionSpec]
	if prevFnSpecJson == "" {
		// is a new deployment or is an old deployment that is missing the annotation
		return true
	}

	prevFnSpec := &faasv1.FunctionSpec{}
	err := json.Unmarshal([]byte(prevFnSpecJson), prevFnSpec)
	if err != nil {
		glog.Errorf("Failed to parse previous function spec: %s", err.Error())
		return true
	}
	prevFn := faasv1.Function{
		Spec: *prevFnSpec,
	}

	if diff := cmp.Diff(prevFn.Spec, function.Spec); diff != "" {
		glog.V(2).Infof("Change detected for %s diff\n%s", function.Name, diff)
		return true
	} else {
		glog.V(3).Infof("No changes detected for %s", function.Name)
	}

	return false
}

func int32p(i int32) *int32 {
	return &i
}

// configureReadOnlyRootFilesystem will create or update the required settings and mounts to ensure
// that the ReadOnlyRootFilesystem setting works as expected, meaning:
// 1. when ReadOnlyRootFilesystem is true, the security context of the container will have ReadOnlyRootFilesystem also
//    marked as true and a new `/tmp` folder mount will be added to the deployment spec
// 2. when ReadOnlyRootFilesystem is false, the security context of the container will also have ReadOnlyRootFilesystem set
//    to false and there will be no mount for the `/tmp` folder
//
// This method is safe for both create and update operations.
func configureReadOnlyRootFilesystem(function *faasv1.Function, deployment *appsv1beta2.Deployment) {
	if deployment.Spec.Template.Spec.Containers[0].SecurityContext != nil {
		deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem = &function.Spec.ReadOnlyRootFilesystem
	} else {
		deployment.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
			ReadOnlyRootFilesystem: &function.Spec.ReadOnlyRootFilesystem,
		}
	}

	existingVolumes := removeVolume("temp", deployment.Spec.Template.Spec.Volumes)
	deployment.Spec.Template.Spec.Volumes = existingVolumes

	existingMounts := removeVolumeMount("temp", deployment.Spec.Template.Spec.Containers[0].VolumeMounts)
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = existingMounts

	if function.Spec.ReadOnlyRootFilesystem {
		deployment.Spec.Template.Spec.Volumes = append(
			existingVolumes,
			corev1.Volume{
				Name: "temp",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		)
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			existingMounts,
			corev1.VolumeMount{Name: "temp", MountPath: "/tmp", ReadOnly: false},
		)
	}
}
