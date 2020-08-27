proxy:
	go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,https://mirrors.aliyun.com/goproxy/,direct   #[翻墙需要]设置代理

build: proxy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/dmp-upload-tool