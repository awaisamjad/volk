package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"volk/config"
	"volk/internal/http"
)

var (
	configPath string
	rootCmd    = &cobra.Command{
		Use:   "volk",
		Short: "Volk is a lightweight HTTP server",
		Long: `Volk is a lightweight HTTP server written in Go, designed to serve static files with minimal configuration.`,
		Run: runServer,
	}
)
var (
	configFlag = flag.String("config", "", "Path to configuration file")
)

func main() {
	flag.Parse()

	// Determine config file path
	configPath := *configFlag
	if configPath == "" {
		// Try standard locations
		possiblePaths := []string{
			"./config/server.toml",
			"./server.toml",
			"/etc/volk/server.toml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	// If no config file found, create default
	if configPath == "" {
		configPath = "./server.toml"
		log.Println("Failed to find config file. Creating one now...")
		err := config.CreateDefaultConfigFile(configPath)
		if err != nil {
			log.Printf("Warning: Failed to create default config file: %v", err)
		} else {
			log.Printf("Created default config file at %s", configPath)
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Configure logging
	setupLogging(cfg.Logging)

	log.Printf("Starting HTTP server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

	// Create file server with configuration
	fileServer := http.NewFileServer(cfg.FileServer)

	// Store file server in a package-level variable or context
	http.DefaultFileServer = fileServer

	fmt.Printf("Listening on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
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
	// Configure log output
	var logOutput *os.File
	var err error

	if logConfig.FilePath != "" {
		// Ensure directory exists
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

	// Set log flags based on format
	if logConfig.Format == "json" {
		// For JSON logging you might want to use a structured logging library
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
}
func handleConnection(conn net.Conn, cfg config.Config) {
	defer conn.Close()

	// ? apply read timeout from config
	if cfg.Server.ReadTimeout > 0 {
		deadline := time.Now().Add(time.Duration(cfg.Server.ReadTimeout) * time.Second)
		conn.SetReadDeadline(deadline)
	}
	reader := bufio.NewReader(conn)

	//? build the request into the correct format string
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

	// Write response back to client
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
