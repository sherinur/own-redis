build:
	gofumpt -l -w .
	go mod tidy
	go build -o own-redis cmd/main.go

clean:
	rm own-redis
	go mod tidy

run:
	gofumpt -l -w .
	go mod tidy
	go build -o own-redis cmd/main.go
	./own-redis
