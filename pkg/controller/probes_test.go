package controller

import (
	"github.com/openfaas/faas-netes/types"
	"testing"
)

func Test_makeProbes_useExec(t *testing.T) {
	readConfig := types.ReadConfig{}
	config := readConfig.Read(types.OsEnv{})

	probes := makeProbes(config)

	if probes.Readiness.Exec == nil {
		t.Errorf("Readiness probe should have had exec handler")
		t.Fail()
	}
	if probes.Liveness.Exec == nil {
		t.Errorf("Liveness probe should have had exec handler")
		t.Fail()
	}
}

func Test_makeProbes_useHTTPProbe(t *testing.T) {
	readConfig := types.ReadConfig{}
	config := readConfig.Read(types.OsEnv{})
	config.HTTPProbe = true

	probes := makeProbes(config)

	if probes.Readiness.HTTPGet == nil {
		t.Errorf("Readiness probe should have had HTTPGet handler")
		t.Fail()
	}
	if probes.Liveness.HTTPGet == nil {
		t.Errorf("Liveness probe should have had HTTPGet handler")
		t.Fail()
	}
}
