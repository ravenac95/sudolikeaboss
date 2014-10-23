package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Name = "sudolikeaboss"
	app.Usage = "use 1password from the terminal with ease"
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "run the sudolikeaboss server for 1password5 workaround",
			Action: func(c *cli.Context) {
				runServer()
			},
		},
		{
			Name:  "getpassword",
			Usage: "Get's the password, like a boss",
			Action: func(c *cli.Context) {
				runGetPassword()
			},
		},
	}

	app.Run(os.Args)
}
