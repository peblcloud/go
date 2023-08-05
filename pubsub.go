package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

func Publish(topic, data string) error {
	res, err := send(&requestArgs{
		method: "POST",
		path:   "publish",
		body: map[string]string{
			topic: data,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during Publish(%s, %s)", topic, data))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status != 0 {
		println(fmt.Sprintf("Exception during Publish(%s, %s)", topic, data))
		println(body.Error)
		return errors.New(body.Error)
	}

	return nil
}

func Subscribe(topic string, cb func(string)) error {
	res, err := send(&requestArgs{
		method: "GET",
		path:   "subscribe",
		query: map[string]string{
			"topic": topic,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during Subscribe(%s)", topic))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode == 200 && body.Status == 0 {
		return nil
	}

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during Subscribe(%s)", topic))
		println(body.Error)
		return errors.New(body.Error)
	}

	if body.Status != 2 {
		println(fmt.Sprintf("Exception during Subscribe(%s)", topic))
		println("received an unrecognized payload, are you perhaps running on an outdated version?")
		return errors.New("unrecognized response")
	}

	numErrors := 0
	for {
		res, _ := send(&requestArgs{
			method: "GET",
			path:   "subscribe_get",
			query: map[string]string{
				"topic": topic,
			},
		})

		body := struct {
			Status int    `json:"status"`
			Error  string `json:"error"`
			Data   string `json:"data"`
		}{}
		json.NewDecoder(res.Body).Decode(&body)

		if body.Status == 0 {
			numErrors = 0
			go cb(body.Data)
		} else {
			if numErrors > 10 {
				println(fmt.Sprintf("encountered too many errors within subscription for %s", topic))
				println("aborting...")
				os.Exit(1)
			}
			println(fmt.Sprintf("encountered error within subscription for %s", topic))
			println("retrying with backoff...")

			time.Sleep(time.Second * (1 << numErrors))
			numErrors += 1
		}
	}
}
