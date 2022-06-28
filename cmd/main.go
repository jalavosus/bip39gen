package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bip39gen",
		Usage: "Generate ethereum addresses which all use separate and randomly-generated seeds, entropy, and mnemonics",
		Commands: []*cli.Command{
			&genCmd,
		},
		DefaultCommand:         genCmdName,
		UseShortOptionHandling: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
