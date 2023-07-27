package pebl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func KVGet(key string) (string, bool, error) {
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	query := url.Values{}
	query.Add("token", token)
	query.Add("key", key)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/kv", kernelURL, kernelPort), nil)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return "", false, errors.New("unable to access kernel")
	}

	var payload [2048]byte
	read, err := res.Body.Read(payload[:])
	if read <= 2 {
		return "", false, errors.New("unable to access kernel")
	}

	if payload[0] == '0' {
		return string(payload[2:read]), true, nil
	}

	if payload[0] == '1' {
		return "", false, nil
	}

	return "", false, errors.New(string(payload[2:read]))
}

func KVSet(key, value string) error {
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	body := url.Values{}
	body.Add(key, value)

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/kv", kernelURL, kernelPort), strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	query := url.Values{}
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("unable to access kernel")
	}

	if res.StatusCode == 200 {
		return nil
	}

	var payload [1024]byte
	read, err := res.Body.Read(payload[:])
	if read > 0 {
		return errors.New(string(payload[:read]))
	}

	return errors.New("unable to access kernel")
}
