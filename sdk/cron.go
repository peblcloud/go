package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	context := os.Getenv("__PEBL_CONTEXT")

	if context == "" {
		kernelURL := os.Getenv("__PEBL_KERNEL_URL")
		kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
		token := os.Getenv("__PEBL_TOKEN")

		query := url.Values{}
		query.Add("token", token)
		query.Add("schedule", schedule)
		query.Add("name", name)

		req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/cron", kernelURL, kernelPort), nil)
		req.URL.RawQuery = query.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			return errors.New(fmt.Sprintf("unable to create cron: %s schedule %s", name, schedule))
		}

		return nil
	}

	if context == name {
		method()
		os.Exit(0)
	}

	return nil
}
