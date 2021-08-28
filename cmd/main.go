package main

import (
	"os"
	"strings"

	"github.com/beeleelee/go-zip/zip"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
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
		log.Fatal(err)
	}
}

var zipCmd = &cli.Command{
	Name:  "zip",
	Usage: "gozip zip [input] [output]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "",
		},
	},
	Action: func(c *cli.Context) error {
		return zip.Zip(c.Args().Get(0), c.Args().Get(1), c.Bool("verbose"))
	},
}

var unzipCmd = &cli.Command{
	Name:  "unzip",
	Usage: "gozip unzip [input] [output]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "",
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "",
		},
		&cli.BoolFlag{
			Name:    "delete-src",
			Aliases: []string{"d"},
			Usage:   "",
		},
	},
	Action: func(c *cli.Context) error {
		source := c.Args().Get(0)
		if !strings.HasSuffix(source, ".zip") {
			return xerrors.New("not a zip file")
		}
		target := c.Args().Get(1)
		if target == "" || target == "./" {
			if curdir, err := os.Getwd(); err != nil {
				return err
			} else {
				target = curdir
			}
		}
		_, err := zip.Unzip(source, target, c.Bool("verbose"), c.Bool("recursive"), c.Bool("delete-src"))
		if err != nil {
			return err
		}
		return nil
	},
}
