package pebl

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Mount(name, path string) error {
	conn, err := rawSend(&requestArgs{
		method: "GET",
		path: "mount",
		query: map[string]string{
			"name": name,
		},
	})

	if err != nil {
		println(fmt.Sprintf("Exeption during Mount(%s, %s)", name, path))
		println(err.Error())
		return errors.New("unable to access the kernel")
	}

	buf := make([]byte, 4096)
	conn.Read(buf[:1])

	if buf[0] != '0' {
		read, _ := conn.Read(buf[:])
		println(fmt.Sprintf("Exeption during Mount(%s, %s)", name, path))
		println(string(buf[:read]))
		return errors.New(string(buf[:read]))
	}

	err = os.MkdirAll(path, 0o777)
	if err != nil {
		println(fmt.Sprintf("Exeption during Mount(%s, %s)", name, path))
		println(fmt.Sprintf("is %s a valid path?", path))
		return errors.New("unable to create folder")
	}

	reader := tar.NewReader(conn)
	for {
		header, err := reader.Next()
		if err != nil {
			break
		}

		switch header.Typeflag {
		case tar.TypeDir:
			path := filepath.Join(path, header.Name)
			os.MkdirAll(path, 0o777)
		case tar.TypeReg | tar.TypeRegA:
			path := filepath.Join(path, header.Name)
			
			f, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				println(fmt.Sprintf("Exeption during Mount(%s, %s)", name, path))
				println(fmt.Sprintf("unable to create %s", header.Name))
				return errors.New("unable to create file")
			}

			io.Copy(f, reader)
		default:
			println("[WARN] unknown file type %d for %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
