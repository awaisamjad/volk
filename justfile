build:
    mkdir -p bin
    GOOS=linux   GOARCH=amd64 go build -o bin/volk_linux_amd64       main.go
    GOOS=darwin  GOARCH=amd64 go build -o bin/volk_darwin_amd64      main.go
    GOOS=windows GOARCH=amd64 go build -o bin/volk_windows_amd64.exe main.go
release:
		gh release create v1.0.0 './bin/volk#Linux'
