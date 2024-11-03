run:
	@go run cmd/main.go

test_lexer:
	@echo "Testing the lexer..."
	@go test -cover ./internal/lexer

test_parser:
	@echo "Testing the parser..."
	@go test -cover ./internal/parser

test: test_lexer test_parser 

coverage:
	@bash scripts/coverage.sh
