package pebl

import (
	"encoding/json"
	"errors"
	"fmt"
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
	res, err := send(&requestArgs{
		method: "GET",
		path:   "mysql",
		query: map[string]string{
			"name": name,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exeption during Mysql(%s)", name))
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
		println(fmt.Sprintf("Exeption during Mysql(%s)", name))
		println(body.Error)
		return nil, errors.New(body.Error)
	}

	return &MysqlConn{
		User:     "root",
		Password: "",
		Addr:     fmt.Sprintf("%s:%d", body.Host, body.Port),
		Host:     body.Host,
		Port:     body.Port,
	}, nil
}
