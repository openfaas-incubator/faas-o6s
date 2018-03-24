package server

import (
	"encoding/json"
	"net/http"

	"github.com/openfaas-incubator/faas-o6s/pkg/version"
	"github.com/openfaas/faas-provider/types"
)

// makeInfoHandler provides the system/info endpoint
func makeInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sha, release := version.GetReleaseInfo()
		info := types.InfoRequest{
			Orchestration: "kubernetes",
			Provider:      "faas-o6s",
			Version: types.ProviderVersion{
				SHA:     sha,
				Release: release,
			},
		}

		infoBytes, _ := json.Marshal(info)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(infoBytes)
	}
}
