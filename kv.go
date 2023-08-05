package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
)

func KVGet(key string) (string, bool, error) {
	res, err := send(&requestArgs{
		method: "GET",
		path:   "kv",
		query: map[string]string{
			"key": key,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during KVGet(%s)", key))
		println("unable to access kernel")
		return "", false, errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
		Data   string `json:"data"`
		Found  bool   `json:"found"`
	}{}
	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during KVGet(%s)", key))
		println(body.Error)
		return "", false, errors.New(body.Error)
	}

	return body.Data, body.Found, nil
}

func KVSet(key, value string) error {
	res, err := send(&requestArgs{
		method: "POST",
		path:   "kv",
		body: map[string]string{
			key: value,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during KVGet(%s)", key))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}
	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during KVGet(%s)", key))
		println(body.Error)
		return errors.New(body.Error)
	}

	return nil
}
