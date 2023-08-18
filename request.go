package pebl

import (
	"fmt"
	"io"
	"net"
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
	request.Header["VERSION"] = []string{version}
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
		os.Setenv("__PEBL_TOKEN", newToken)
	}

	return res, err
}

func rawSend(args *requestArgs) (net.Conn, error) {
	host, port, token := checkEnv()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	var ps string
	if len(args.query) > 0 {
		values := url.Values{}
		for k, v := range args.query {
			values.Add(k, v)
		}
		ps = fmt.Sprintf("/%s?%s", args.path, values.Encode())
	} else {
		ps = fmt.Sprintf("/%s", args.path)
	}

	conn.Write([]byte(fmt.Sprintf("%s %s HTTP/1.1\r\n", args.method, ps)))
	conn.Write([]byte(fmt.Sprintf("Host: %s:%s\r\n", host, port)))
	conn.Write([]byte(fmt.Sprintf("TOKEN: %s\r\n", token)))
	conn.Write([]byte(fmt.Sprintf("VERSION: %s\r\n", version)))

	if len(args.body) > 0 {
		conn.Write([]byte("Content-Type: application/x-www-form-urlencoded\r\n"))
		values := url.Values{}
		for k, v := range args.body {
			values.Add(k, v)
		}
		payload := []byte(values.Encode())
		length := len(payload)
		conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", length)))
		conn.Write(payload)
	} else {
		conn.Write([]byte("Content-Length: 0\r\n\r\n"))
	}

	return conn, nil
}
