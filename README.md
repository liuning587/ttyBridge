# ttyBridge TTY转TCP代理工具

ttyBridge

## Go mod

- Init: go mod init ccoRouter
- download: go mod download
- vendor: go mod vendor
- tidy: go mod tidy

## Depend

- go get github.com/jacobsa/go-serial
- go get golang.org/x/sys

## linux

1. CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w"
2. upx -9 ttyBridge

## windows

1. SET GOARCH=arm
2. SET GOOS=linux
3. SET GOARM=7
4. SET CGO_ENABLED=1
5. SET CC=arm-linux-gnueabihf-gcc
6. go build / go build -ldflags "-s -w"
7. upx -9 ttyBridge
