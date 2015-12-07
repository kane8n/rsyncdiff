BUILDDATE=`date -u +.%Y%m%d%.H%M%S`
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.builddate ${BUILDDATE}" -o="build/linux_amd64/rsyncdiff" rsyncdiff.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.builddate ${BUILDDATE}" -o="build/darwin_amd64/rsyncdiff" rsyncdiff.go
