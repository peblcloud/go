package sdk

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

func makeRequest(path, mode string) (net.Conn, error) {
	kernelURL := os.Getenv("__PEBL_KERNEL_URL")
	kernelPort := os.Getenv("__PEBL_KERNEL_PORT")
	token := os.Getenv("__PEBL_TOKEN")

	query := url.Values{}
	query.Add("token", token)
	query.Add("path", path)
	query.Add("mode", mode)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/open", kernelURL, kernelPort), nil)
	req.URL.RawQuery = query.Encode()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", kernelURL, kernelPort))
	if err != nil {
		return nil, errors.New("unable to reach the kernel")
	}

	conn.Write([]byte(fmt.Sprintf("GET /open?%s HTTP/1.1\r\n", query.Encode())))
	conn.Write([]byte(fmt.Sprintf("Host: %s:%s\r\n", kernelURL, kernelPort)))
	conn.Write([]byte("Content-Length: 0\r\n\r\n"))

	var payload [1]byte
	read, _ := conn.Read(payload[:])
	if read != 1 || payload[0] != '0' {
		return nil, errors.New("error requesting the kernel")
	}

	return conn, nil
}

// Write returns an object with io.WriteCloser interface.
// Close() must be called before the contents of the Write()
// are made available at the provided path.
//
// Interleaving Read and Write on the same path may result in
// unexpected behavior.
func Write(path string) (io.WriteCloser, error) {
	conn, err := makeRequest(path, "w")
	if err != nil {
		return nil, err
	}
	return conn.(io.WriteCloser), nil
}

// Read returns an object with io.ReadCloser interface.
//
// Interleaving Read and Write on the same path may result in
// unexpected behavior.
func Read(path string) (io.ReadCloser, error) {
	conn, err := makeRequest(path, "r")
	if err != nil {
		return nil, err
	}
	return conn.(io.ReadCloser), nil
}
