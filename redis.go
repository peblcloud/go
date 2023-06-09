package pebl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	query := url.Values{}
	query.Add("token", token)
	query.Add("name", name)

	redisReq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/redis", kernelURL, kernelPort), nil)
	redisReq.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(redisReq)
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unable to create redis instance %s", name))
	}

	var payload [1024]byte
	read, err := res.Body.Read(payload[:])
	if read <= 2 {
		return nil, errors.New(fmt.Sprintf("unable to read payload from kernel: %s", err.Error()))
	}

	if payload[0] != '0' {
		return nil, errors.New(fmt.Sprintf("unable to create redis instance %s", name))
	}

	addr := string(payload[2:read])
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return nil, errors.New("unable to parse payload from kernel")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("unable to parse payload from kernel")
	}

	return &RedisConn{
		Addr: addr,
		Host: parts[0],
		Port: port,
	}, nil
}
