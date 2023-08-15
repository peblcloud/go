package pebl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var version string = "0.1.0"

type requestArgs struct {
	method string
	path   string
	query  map[string]string
	body   map[string]string
}

func checkEnv() (string, string, string) {
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	if kernelURL == "" || kernelPort == "" || token == "" {
		println("no pebl context was detected, aborting!")
		println("pebl programs must be run within a pebl environment, either with")
		println("a local pebl cluster or in the pebl cloud environment.")
		println("")
		println("for more information visit: https://docs.pebl.io/issues")
		println("")
		os.Exit(1)
	}

	return kernelURL, kernelPort, token
}

func send(args *requestArgs) (*http.Response, error) {
	host, port, token := checkEnv()

	addr := fmt.Sprintf("http://%s:%s/%s", host, port, args.path)

	var body io.Reader = nil
	if len(args.body) > 0 {
		values := url.Values{}
		for k, v := range args.body {
			values.Add(k, v)
		}
		body = strings.NewReader(values.Encode())
	}

	request, _ := http.NewRequest(args.method, addr, body)
	request.Header["TOKEN"] = []string{token}
	if len(args.body) > 0 {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	values := url.Values{}
	for k, v := range args.query {
		values.Add(k, v)
	}
	request.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(request)

	if newToken := res.Header.Get("Token"); newToken != "" {
		println(" -- new token")
		os.Setenv("__PEBL_TOKEN", newToken)
	}

	return res, err
}
