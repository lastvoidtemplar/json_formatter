package main

import "github.com/lastvoidtemplar/json_formatter/internal/server"

func main() {
	server.New().Run(3000)
}
