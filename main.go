package main

import "mcp-gateway/src/module"

func main() {
	app := module.NewApp()
	app.Run()
}
