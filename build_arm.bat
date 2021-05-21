SET GOARCH=arm
SET GOOS=linux
SET GOARM=7
SET CGO_ENABLED=1
SET CC=arm-linux-gnueabihf-gcc
go build -ldflags "-s -w"
upx -9 ttyBridge