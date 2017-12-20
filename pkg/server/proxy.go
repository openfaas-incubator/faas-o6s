package server

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas/gateway/requests"
)

// makeProxy creates a proxy for HTTP web requests which can be routed to a function.
func makeProxy(functionNamespace string) http.HandlerFunc {
	proxyClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body != nil {
			defer r.Body.Close()
		}

		switch r.Method {
		case http.MethodGet,
			http.MethodPost:

			vars := mux.Vars(r)
			service := vars["name"]

			defer func(when time.Time) {
				seconds := time.Since(when).Seconds()
				glog.Infof("%s took %f seconds", service, seconds)
			}(time.Now())

			var addr string

			entries, lookupErr := net.LookupIP(fmt.Sprintf("%s.%s", service, functionNamespace))
			if lookupErr == nil && len(entries) > 0 {
				index := randomInt(0, len(entries))
				addr = entries[index].String()
			}

			forwardReq := requests.NewForwardRequest(r.Method, *r.URL)

			url := forwardReq.ToURL(addr, 8080)

			request, _ := http.NewRequest(r.Method, url, r.Body)

			copyHeaders(&request.Header, &r.Header)

			defer request.Body.Close()

			response, err := proxyClient.Do(request)

			if err != nil {
				glog.Errorf("%s error: %s", service, err.Error())
				writeHead(service, http.StatusInternalServerError, w)
				buf := bytes.NewBufferString("Can't reach service: " + service)
				w.Write(buf.Bytes())
				return
			}

			clientHeader := w.Header()
			copyHeaders(&clientHeader, &response.Header)

			writeHead(service, http.StatusOK, w)
			io.Copy(w, response.Body)
		}
	}
}

func writeHead(service string, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
}

func copyHeaders(destination *http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(*destination)[k] = vClone
	}
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
