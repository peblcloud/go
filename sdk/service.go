// Package sdk provides go bindings for pebl's system calls.
//
// For more detailed guides on utilizing pebl, check out [the docs]
//
// [the docs]: https://docs.pebl.io/
package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Service takes an http.Handler and creates a serving
// resource at the endpoint.
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

	http.ListenAndServe(":8000", mux)
	return nil
}
