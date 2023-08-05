package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Cron creates a scheduled task.
//
// The name identifies the scheduled task, and subsequent
// uses of the same name will overwrite the same logical
// task.
//
// The schedule must be valid cron schedule with 5 fields
// separated by spaces. It also accepts @hourly and @daily,
// which are shorthand for "0 * * * *" and "0 0 * * *" respectively.
func Cron(name, schedule string, method func()) error {
	res, err := send(&requestArgs{
		method: "GET",
		path:   "cron",
		query: map[string]string{
			"name":     name,
			"schedule": schedule,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exception during Cron(%s, %s)", name, schedule))
		println("unable to access kernel")
		return errors.New("unable to access kernel")
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{}
	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status == 1 {
		println(fmt.Sprintf("Exception during Cron(%s, %s)", name, schedule))
		println(body.Error)
		return errors.New(body.Error)
	}

	if body.Status == 0 {
		return nil
	}

	if body.Status != 2 {
		println(fmt.Sprintf("Exception during Cron(%s, %s)", name, schedule))
		println("received an unrecognized payload, are you perhaps running on an outdated version?")
		return errors.New("unrecognized response")
	}

	method()
	os.Exit(0)
	return nil
}
