// Package sdk provides go bindings for pebl's system calls.
//
// For more detailed guides on utilizing pebl, check out [the docs]
//
// [the docs]: https://docs.pebl.io/
package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Service takes an http.Handler and creates a serving
// resource at the endpoint. This service will be reachable
// from outside the cluster from the provided `endpoint` argument.
//
// The endpoint must be a valid domain owned by the
// user.
func Service(app http.Handler, endpoint string) error {
	res, err := send(&requestArgs{
		method: "GET",
		path:   "service",
		query: map[string]string{
			"endpoint": endpoint,
			"internal": "0",
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during Service(app, %s)", endpoint))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during Service(app, %s)", endpoint))
		println(body.Error)
		return errors.New(body.Error)
	}

	if body.Status == 0 {
		return nil
	}

	if body.Status != 2 {
		println(fmt.Sprintf("Exception during Service(app, %s)", endpoint))
		println("received an unrecognized payload, are you perhaps running on an outdated version?")
		return errors.New("unrecognized response")
	}

	s := http.Server{
		Addr:    ":80",
		Handler: app,
	}
	log.Fatal(s.ListenAndServe())
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
	res, err := send(&requestArgs{
		method: "GET",
		path:   "service",
		query: map[string]string{
			"endpoint": endpoint,
			"internal": "1",
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during InternalService(app, %s)", endpoint))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during InternalService(app, %s)", endpoint))
		println(body.Error)
		return errors.New(body.Error)
	}

	if body.Status == 0 {
		return nil
	}

	if body.Status != 2 {
		println(fmt.Sprintf("Exception during InternalService(app, %s)", endpoint))
		println("received an unrecognized payload, are you perhaps running on an outdated version?")
		return errors.New("unrecognized response")
	}

	s := http.Server{
		Addr:    ":80",
		Handler: app,
	}
	log.Fatal(s.ListenAndServe())
	return nil
}
