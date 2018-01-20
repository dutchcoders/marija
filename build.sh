LDFLAGS="$(go run -exec ~/.gopath/bin/sign-wrapper.sh scripts/gen-ldflags.go)"
env GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ./bin/marija-linux-amd64

