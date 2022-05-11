config:
	@cp config.example.json config.json

test:
	go test -v -cover ./...