package specutils

import (
	"strconv"

	"github.com/golang/glog"
)

const (
	defaultMinReplicas   = 1
	defaultMaxReplicas   = 100
	defaultScalingFactor = 20
	defaultScaleToZero   = false
	// LabelMinReplicas is used to set the min replica count a function can have, default 1.
	LabelMinReplicas = "com.openfaas.scale.min"
	// LabelMaxReplicas is used to set the max replica count a function can have, default 100.
	LabelMaxReplicas = "com.openfaas.scale.max"
	// LabelScalingFactor is used to precentage step size during scaling, default 20.
	// When zero, it disables OpenFaaS auto-scaling.
	LabelScalingFactor = "com.openfaas.scale.factor"
	// LabelScaleToZero is used to enable OpenFaaS to scale the function to zero, i.e. idle the function,
	// when there is no traffic.
	LabelScaleToZero = "com.openfaas.scale.zero"
)

// GetMinReplicaCount parses the function min allowed replicas value form its labels
func GetMinReplicaCount(labels *map[string]string) int32 {
	value := getLabelInt(labels, LabelMinReplicas, 0)
	if value > 0 {
		return int32(value)
	}

	return defaultMinReplicas
}

// GetMaxReplicaCount parses the function max allowed replicas value form its labels
func GetMaxReplicaCount(labels *map[string]string) int32 {
	value := getLabelInt(labels, LabelMaxReplicas, 0)

	if value > 0 {
		return int32(value)
	}

	return defaultMaxReplicas
}

// GetScalingFactor parses the function scaling factor label. The value is between 0 and 100,
// the default value is 20.  The 0 value is used to disable OpenFaaS scaling of the function.
func GetScalingFactor(labels *map[string]string) int {
	value := getLabelInt(labels, LabelScalingFactor, defaultScalingFactor)

	if value < 0 {
		return 0
	}

	if value > 100 {
		return 100
	}

	return value
}

// getLabelInt parses the integer value of the supplied label.  If value does not exist or is
// not a valid integer, notFound is returned
func getLabelInt(labels *map[string]string, name string, notFound int) int {
	if labels == nil {
		return notFound
	}

	valueStr, exists := (*labels)[name]
	if !exists {
		return notFound
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		glog.Error(err)
		return notFound
	}
	return value
}

// GetScaleToZero parses the scale to zero label. Returns false by default.
func GetScaleToZero(labels *map[string]string) bool {
	return getLabelBool(labels, LabelScaleToZero, defaultScaleToZero)
}

// getLabelBool parses the boolean value of the supplied label. If the label is not found or
// can not be parsed, notFound is returned.
func getLabelBool(labels *map[string]string, name string, notFound bool) bool {
	if labels == nil {
		return notFound
	}

	valueStr, exists := (*labels)[name]
	if !exists {
		return notFound
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		glog.Error(err)
		return notFound
	}
	return value
}
