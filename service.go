// Package sdk provides go bindings for pebl's system calls.
//
// For more detailed guides on utilizing pebl, check out [the docs]
//
// [the docs]: https://docs.pebl.io/
package pebl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Service takes an http.Handler and creates a serving
// resource at the endpoint. This service will be reachable
// from outside the cluster from the provided `endpoint` argument.
//
// The endpoint must be a valid domain owned by the
// user.
func Service(app http.Handler, endpoint string) error {
	context := os.Getenv("__PEBL_CONTEXT")

	if context == "" {
		kernelURL := os.Getenv("__PEBL_KERNEL_URL")
		kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
		token := os.Getenv("__PEBL_TOKEN")

		serviceReq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/service", kernelURL, kernelPort), nil)
		query := url.Values{}
		query.Add("token", token)
		query.Add("endpoint", endpoint)
		serviceReq.URL.RawQuery = query.Encode()

		res, err := http.DefaultClient.Do(serviceReq)
		if err != nil {
			return errors.New("unable to create service")
		}
		if res.StatusCode == 200 {
			return nil
		} else {
			return errors.New("unable to create service")
		}
	}

	if context != endpoint {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.ServeHTTP(w, r)
	})

	http.ListenAndServe(":80", mux)
	return nil
}

// InternalService takes an http.Handler and creates a
// private service at the endpoint. A private service can
// only be reached by other workloads in the cluster using
// the provided `endpoint` argument.
//
// Unlike `Service`, the `endpoint` used here can be any valid domain.
// Though in general we recommend a scheme to easily identify internal
// vs. external services, something like `foo.local` or `bar.private`.
func InternalService(app http.Handler, endpoint string) error {
	context := os.Getenv("__PEBL_CONTEXT")

	if context == "" {
		kernelURL := os.Getenv("__PEBL_KERNEL_URL")
		kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
		token := os.Getenv("__PEBL_TOKEN")

		serviceReq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/service", kernelURL, kernelPort), nil)
		query := url.Values{}
		query.Add("token", token)
		query.Add("endpoint", endpoint)
		query.Add("internal", "1")
		serviceReq.URL.RawQuery = query.Encode()

		res, err := http.DefaultClient.Do(serviceReq)
		if err != nil {
			println(err.Error())
			return errors.New("unable to create service")
		}
		if res.StatusCode == 200 {
			return nil
		} else {
			println(res.Status)
			return errors.New("unable to create service")
		}
	}

	if context != endpoint {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.ServeHTTP(w, r)
	})

	http.ListenAndServe(":80", mux)
	return nil
}
