test:
	go vet ./...
	go test -v ./... -count=1
