package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/awaisamjad/volk/config"
	"github.com/awaisamjad/volk/internal/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve files over HTTP",
	Long:  `The serve command starts an HTTP server that serves files from the directory.`,
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	setupLogging(cfg.Logging)

	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.NewFileServer(cfg.FileServer)
	http.DefaultFileServer = fileServer

	fmt.Printf("Listening on localhost:%d\n", cfg.Server.Port)
	fmt.Printf("Serving files from: %s\n", cfg.FileServer.DocumentRoot)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection ", err)
		}
		go handleConnection(conn, cfg)
	}

}

func setupLogging(logConfig config.LogConfig) {
	var logOutput *os.File
	var err error

	if logConfig.FilePath != "" {
		dir := filepath.Dir(logConfig.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Could not create log directory: %v", err)
		}

		logOutput, err = os.OpenFile(logConfig.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: Could not open log file: %v", err)
		} else {
			log.SetOutput(logOutput)
		}
	}
	switch logConfig.Format {
	case "plain":
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	case "verbose":
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}

}

func handleConnection(conn net.Conn, cfg config.Config) {
	defer conn.Close()

	if cfg.Server.ReadTimeout > 0 {
		deadline := time.Now().Add(time.Duration(cfg.Server.ReadTimeout) * time.Second)
		conn.SetReadDeadline(deadline)
	}

	reader := bufio.NewReader(conn)
	var requestBuilder strings.Builder
	startLine, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Error reading start line: %v", err)
		return
	}
	requestBuilder.WriteString(startLine)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading header line: %v", err)
			return
		}

		requestBuilder.WriteString(line)

		if line == "\r\n" || line == "\n" {
			break
		}
	}

	req, err := http.NewRequest(requestBuilder.String())
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request"))
		return
	}

	resp := req.Response()

	_, err = conn.Write([]byte(resp.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}

	if cfg.Logging.AccessLogs {
		log.Printf("Access: %s %s %s - %d %s",
			req.StartLine.Method,
			req.StartLine.RequestTarget,
			req.StartLine.Protocol,
			resp.StartLine.StatusCode,
			resp.StartLine.StatusText)
	}
}
