dotz_ver = 0.0.0
build: main.go
	GOOS=darwin GOARCH=amd64 go build -v -ldflags "-X main.dotzVersion=${dotz_ver}" -o ./dist/dotz-amd64 .
	GOOS=darwin GOARCH=arm64 go build -v -ldflags "-X main.dotzVersion=${dotz_ver}" -o ./dist/dotz-arm64 .
clean:
	rm ./dist/*
