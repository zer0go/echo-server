package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

var Version = "development"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RequestURI)
		host := request.URL.Query().Get("host")
		port := request.URL.Query().Get("port")

		if host == "" || port == "" {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Host and port are required"))
			return
		}

		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("Connecting error: %s", err)))
		}
		if conn != nil {
			defer conn.Close()
			writer.Write([]byte(fmt.Sprintf("Opened: %s", net.JoinHostPort(host, port))))
		}
	})
	mux.HandleFunc("/env", func(writer http.ResponseWriter, request *http.Request) {
		envs := os.Environ()
		fmt.Println(envs)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo Server (%s) is running on port http://0.0.0.0:%s", Version, port)
	err := http.ListenAndServe("0.0.0.0:"+port, mux)
	if err != nil {
		panic(err)
	}
}
