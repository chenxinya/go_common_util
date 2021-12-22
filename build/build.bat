SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
e:
cd ../
go build -o name cmd/main.go