run:
	@go run cmd/main.go

test_lexer:
	@go test -cover ./internal/lexer