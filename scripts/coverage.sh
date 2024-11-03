# in first argument($1) can be specified the output directory of coverage.html
out=${1:-"."}

go test -coverprofile=coverage.out ./internal/lexer

go test -coverprofile=tmp.out ./internal/parser
tail -n +2 tmp.out | grep -v "parser_error" >> coverage.out
rm tmp.out

go tool cover -html=coverage.out -o ${out%%/}/coverage.html
rm coverage.out