package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Name = "sudolikeaboss"
	app.Version = "0.2.0"
	app.Usage = "use 1password from the terminal with ease"
	app.Action = func(c *cli.Context) {
		runSudolikeaboss()
	}

	app.Run(os.Args)
}
