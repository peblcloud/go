package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
)

// RedisConn provides the Host and Port of the target redis instance
type RedisConn struct {
	Addr string
	Host string
	Port int
}

// Redis returns connection information for the provided redis instance
// with the name `name`
func Redis(name string) (*RedisConn, error) {
	res, err := send(&requestArgs{
		method: "GET",
		path:   "redis",
		query: map[string]string{
			"name": name,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exeption during Redis(%s)", name))
		println("unable to access the kernel")
		return nil, err
	}

	body := struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode != 200 || body.Status != 0 {
		println(fmt.Sprintf("Exeption during Redis(%s)", name))
		println(body.Error)
		return nil, errors.New(body.Error)
	}

	return &RedisConn{
		Addr:     fmt.Sprintf("%s:%d", body.Host, body.Port),
		Host:     body.Host,
		Port:     body.Port,
	}, nil

}
