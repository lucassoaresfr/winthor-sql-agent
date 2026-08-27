package main

import "github.com/lucassoaresfr/winthor-sql-agent.git/bootstrap"

func main() {
	app := bootstrap.NewApp()
	app.Run()
}
