package main

import (
	"fmt"
	"os"

	"github.com/karlmutch/aws-ec2-price/pkg/rest"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ec2-instance-price"
	app.Version = "0.0.3"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port",
			EnvVar: "PORT",
			Value:  8080,
		},
	}
	app.Action = func(c *cli.Context) error {
		port := fmt.Sprintf(":%d", c.Int("port"))

		r := rest.GetRouter()
		r.Run(port)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
