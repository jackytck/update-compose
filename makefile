publish: *.go
	env GOOS=linux GOARCH=amd64 go build -o update-compose

lint:
	golint ./...
	golangci-lint run

clean:
	rm update-compose
