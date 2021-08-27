package main

import (
	"fmt"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/it-fm/gozip"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("gozip")

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		zipCmd,
		unzipCmd,
	}

	app := &cli.App{
		Name:     "gozip",
		Flags:    []cli.Flag{},
		Commands: local,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var zipCmd = &cli.Command{
	Name:  "zip",
	Usage: "gozip zip [input] [output]",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		return gozip.Zip(c.Args().Get(0), c.Args().Get(1))
	},
}

var unzipCmd = &cli.Command{
	Name:  "unzip",
	Usage: "gozip unzip [input] [output]",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		lists, err := gozip.Unzip(c.Args().Get(0), c.Args().Get(1))
		if err != nil {
			return err
		}
		for _, li := range lists {
			fmt.Println(li)
		}
		return nil
	},
}
