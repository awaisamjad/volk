package cmd

import (
	"bufio"
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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve files over HTTP",
	Long:  `The serve command starts an HTTP server that serves files from a specified directory.`,
	Run:   runServer,
}

var configPath string
var createConfig bool

func init() {
	serveCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	serveCmd.PersistentFlags().BoolVarP(&createConfig, "createConfig", "C", false, "Create a default config in the root dir")
}

var ServerConfigPossiblePaths = []string{
	"./server.toml",
	"./config/server.toml",
	"/etc/volk/server.toml",
}

func runServer(cmd *cobra.Command, args []string) {

	if configPath == "" {
		found := false
		for _, path := range ServerConfigPossiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				found = true
				break
			}
		}
		if !found {
			if createConfig {
				configPath = "./server.toml"
				log.Println("No config file found in default locations. Creating one now...")
				err := config.CreateDefaultConfigFile(configPath)
				if err != nil {
					log.Fatalf("Failed to create default config file: %v", err)
				} else {
					log.Printf("Created default config file at %s", configPath)
				}
			} else {
				var paths []string
				for _, path := range ServerConfigPossiblePaths {
					paths = append(paths, path)
				}
				log.Fatalf(`Can't locate server file. Provide the path or create one in the following areas : %v, or use --createConfig|-C to create a default one.
				Example: volk serve -C	
				`, paths)
			}
		}
	} else {
		if _, err := os.Stat(configPath); err != nil {
			log.Fatalf("Config file specified but not found: %s", configPath)
		}
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	setupLogging(cfg.Logging)

	log.Printf("Starting HTTP server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.NewFileServer(cfg.FileServer)
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
	switch logConfig.Format{
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
