package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var Version = "development"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		resp := make(map[string]any)
		resp["headers"] = request.Header
		resp["queryParams"] = request.URL.Query()
		b, _ := json.Marshal(resp)

		dump, _ := httputil.DumpRequest(request, true)
		fmt.Printf("%q", dump)

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(b)
	})
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
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
			_, _ = writer.Write([]byte(fmt.Sprintf("Connecting error: %s", err)))
		}
		if conn != nil {
			defer func(conn net.Conn) {
				_ = conn.Close()
			}(conn)
			_, _ = writer.Write([]byte(fmt.Sprintf("Opened: %s", net.JoinHostPort(host, port))))
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

	fmt.Printf("Echo Server (%s) is running on port http://0.0.0.0:%s\n", Version, port)
	err := http.ListenAndServe("0.0.0.0:"+port, mux)
	if err != nil {
		panic(err)
	}
}
