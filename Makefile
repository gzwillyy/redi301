# BINARY FILE
PROJECT="redi301"
# START FILE PATH
MAIN_PATH="main.go"

build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/${PROJECT} ${MAIN_PATH}
