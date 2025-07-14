package http

import (
	"fmt"
	"log"

	// "server/internal/
	"mime"
	"os"
	"path"
	"path/filepath"
	"volk/config"
)

// FileServer handles serving files
type FileServer struct {
	Config config.FileServerConfig
}

func NewFileServer(config config.FileServerConfig) *FileServer {
	return &FileServer{
		Config: config,
	}
}

// ServeFile handles file serving based on an request
func (fs *FileServer) ServeFile(req *Request) Response {
	if req.GetMethod() != GET {
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   req.StartLine.Protocol,
				StatusCode: 405,
				StatusText: StatusCodeMap[405],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
				{Name: "Allow", Value: "GET"},
			},
			Body: "405 Method Not Allowed: Only GET is supported for file serving",
		}
	}

	urlPath := req.GetRequestTarget()
	err := req.ValidatePath()
	if err != nil {
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   req.StartLine.Protocol,
				StatusCode: 400,
				StatusText: StatusCodeMap[400],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "400 Bad Request: Invalid path",
		}
	}

	cleanPath := path.Clean(urlPath)
	filePath := filepath.Join(fs.Config.DocumentRoot, cleanPath[1:])
	// filePath = cleanPath[1:] // Remove leading slash
	//? check if it exists and is a directory
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(err)
			return Response{
				StartLine: ResponseStartLine{
					Protocol:   req.StartLine.Protocol,
					StatusCode: 404,
					StatusText: StatusCodeMap[404],
				},
				Headers: []Header{
					{Name: "Content-Type", Value: "text/plain"},
				},
				Body: "404 Not Found",
			}
		}
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   req.StartLine.Protocol,
				StatusCode: 500,
				StatusText: StatusCodeMap[500],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "500 Internal Server Error",
		}
	}

	//? if the path given is a dir then check if it has an index.html which will be served
	if fileInfo.IsDir() {
		filePath = filepath.Join(filePath, fs.Config.DefaultFile)
		log.Println("isDir", filePath)
		_, err := os.Stat(filePath)
		if err != nil {
			log.Println(err)
			return Response{
				StartLine: ResponseStartLine{
					Protocol:   req.StartLine.Protocol,
					StatusCode: 403,
					StatusText: StatusCodeMap[403],
				},
				Headers: []Header{
					{Name: "Content-Type", Value: "text/plain"},
				},
				Body: "403 Forbidden: Directory listing not allowed",
			}
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   req.StartLine.Protocol,
				StatusCode: 500,
				StatusText: StatusCodeMap[500],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "500 Internal Server Error",
		}
	}

	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return Response{
		StartLine: ResponseStartLine{
			Protocol:   req.StartLine.Protocol,
			StatusCode: 200,
			StatusText: StatusCodeMap[200],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: contentType},
			{Name: "Content-Length", Value: fmt.Sprintf("%d", len(content))},
		},
		Body: string(content),
	}
}
