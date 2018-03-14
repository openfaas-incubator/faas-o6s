package controller

import (
	"fmt"

	faasv1alpha1 "github.com/openfaas-incubator/faas-o6s/pkg/apis/o6sio/v1alpha1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
)

// UpdateSecrets will update the Deployment spec to include secrets that have been deployed
// in the kubernetes cluster.  For each requested secret, we inspect the type and add it to the
// deployment spec as appropriate: secrets with type `SecretTypeDockercfg` are added as ImagePullSecrets
// all other secrets are mounted as files in the deployments containers.
func UpdateSecrets(function *faasv1alpha1.Function, deployment *appsv1beta2.Deployment, existingSecrets map[string]*corev1.Secret) error {

	// Add / reference pre-existing secrets within Kubernetes
	secretVolumeProjections := []corev1.VolumeProjection{}
	for _, secretName := range function.Spec.Secrets {
		deployedSecret, ok := existingSecrets[secretName]
		if !ok {
			return fmt.Errorf("required secret '%s' was not found in the cluster", secretName)
		}

		if deployedSecret.Type == corev1.SecretTypeDockercfg {
			deployment.Spec.Template.Spec.ImagePullSecrets = append(
				deployment.Spec.Template.Spec.ImagePullSecrets,
				corev1.LocalObjectReference{
					Name: secretName,
				},
			)
		} else {
			// projectSecrets.VolumeSource.Sources = newProjections
			projectedPaths := []corev1.KeyToPath{}
			for secretKey := range deployedSecret.Data {
				projectedPaths = append(projectedPaths, corev1.KeyToPath{Key: secretKey, Path: secretKey})
			}

			projection := &corev1.SecretProjection{Items: projectedPaths}
			projection.Name = secretName
			secretProjection := corev1.VolumeProjection{
				Secret: projection,
			}
			secretVolumeProjections = append(secretVolumeProjections, secretProjection)

		}
	}

	if len(secretVolumeProjections) > 0 {
		volumeName := fmt.Sprintf("%s-projected-secrets", function.Name)
		projectedSecrets := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Projected: &corev1.ProjectedVolumeSource{
					Sources: secretVolumeProjections,
				},
			},
		}
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, projectedSecrets)

		// add mount secret as a file
		updatedContainers := []corev1.Container{}
		for _, container := range deployment.Spec.Template.Spec.Containers {
			mount := corev1.VolumeMount{
				Name:      volumeName,
				ReadOnly:  true,
				MountPath: "/run/secrets",
			}
			container.VolumeMounts = append(container.VolumeMounts, mount)
			updatedContainers = append(updatedContainers, container)
		}

		deployment.Spec.Template.Spec.Containers = updatedContainers
	}

	return nil
}
