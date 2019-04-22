package controller

import (
	"github.com/openfaas/faas-netes/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
	"path/filepath"
)

type FunctionProbes struct {
	Liveness  *corev1.Probe
	Readiness *corev1.Probe
}

// makeProbes returns liveness and readiness configured with exec if
// the env var http_probe is false
func makeProbes(config types.BootstrapConfig) *FunctionProbes {
	var handler corev1.Handler

	if config.HTTPProbe {
		handler = corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/_/health",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(functionPort),
				},
			},
		}
	} else {
		path := filepath.Join(os.TempDir(), ".lock")
		handler = corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{"cat", path},
			},
		}
	}

	probes := FunctionProbes{}
	probes.Readiness = &corev1.Probe{
		Handler:             handler,
		InitialDelaySeconds: int32(config.ReadinessProbeInitialDelaySeconds),
		TimeoutSeconds:      int32(config.ReadinessProbeTimeoutSeconds),
		PeriodSeconds:       int32(config.ReadinessProbePeriodSeconds),
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	probes.Liveness = &corev1.Probe{
		Handler:             handler,
		InitialDelaySeconds: int32(config.LivenessProbeInitialDelaySeconds),
		TimeoutSeconds:      int32(config.LivenessProbeTimeoutSeconds),
		PeriodSeconds:       int32(config.LivenessProbePeriodSeconds),
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	return &probes
}
