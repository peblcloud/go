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

// MysqlConn provides the connection information for a mysql instance
type MysqlConn struct {
	User     string
	Password string
	Addr     string
	Host     string
	Port     int
}

// Mysql returns connection information for the provided mysql instance
// with the name `name`
func Mysql(name string) (*MysqlConn, error) {
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	query := url.Values{}
	query.Add("token", token)
	query.Add("name", name)

	mysqlReq, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/mysql", kernelURL, kernelPort), nil)
	mysqlReq.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(mysqlReq)
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unable to create mysql instance %s", name))
	}

	var payload [1024]byte
	read, err := res.Body.Read(payload[:])
	if read <= 2 {
		return nil, errors.New(fmt.Sprintf("unable to read payload from kernel: %s", err.Error()))
	}

	if payload[0] != '0' {
		return nil, errors.New(fmt.Sprintf("unable to create mysql instance %s", name))
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

	return &MysqlConn{
		User: "root",
		Password: "",
		Addr: addr,
		Host: parts[0],
		Port: port,
	}, nil
}
