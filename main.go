package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

type Response struct {
	StatCode uint        `json:"stat_code"`
	StatMsg  string      `json:"stat_msg"`
	Data     interface{} `json:"data,omitempty"`
}

func main() {
	LoadEnvConfig()
	s, n := Server()

	s.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		res := Response{
			StatCode: http.StatusOK,
			StatMsg:  "Building sso and cloud storage using keyloack and minio",
		}

		json.NewEncoder(rw).Encode(&res)
	})

	err := http.Serve(n, s)
	if err != nil {
		log.Fatalf("Server listening error %v", err)
	} else {
		log.Printf("Server listening on port: %s", os.Getenv("GO_PORT"))
	}
}

func Server() (*http.ServeMux, net.Listener) {
	var lc = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			err := c.Control(func(fd uintptr) {
				err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				if err != nil {
					log.Fatal(err)
				}
			})
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("app:%s", os.Getenv("GO_PORT")))
	if err != nil {
		log.Fatalf("Net listening error %v", err)
	}

	router := http.NewServeMux()

	return router, ln
}

func LoadEnvConfig() {
	_, bool := os.LookupEnv("GO_PORT")
	if !bool {
		os.Setenv("GO_PORT", "3000")
	}
}
