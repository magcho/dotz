dotz_ver = 0.0.0
build: main.go
	go build -v -ldflags "-X main.dotzVersion=${dotz_ver}" .
clean:
	rm dotz
