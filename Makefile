run:
	@go run cmd/main.go

test_lexer:
	@go test -cover ./internal/lexer

test_parser:
	@go test -cover ./internal/parser

test test_lexer test_parser: 