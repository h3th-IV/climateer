package main

import (
	"fmt"
	"os"

	"github.com/h3th-IV/climateer/pkg/command"
	"github.com/urfave/cli/v2"
)

func main() {
	climateer := &cli.App{
		Name:  "Climate Data Repository",
		Usage: "The is a repository for climate Data from Africa and around the world",
		Commands: []*cli.Command{
			command.StartCommand(),
		},
		Version: "v.0.1.5",
		Authors: []*cli.Author{
			{
				Name:  "Funmilola Cole",
				Email: "Deborahcolex10@gmail.com",
			},
			{
				Name:  "Suzzy Niniola",
				Email: "...........",
			},
		},
	}
	if err := climateer.Run(os.Args); err != nil {
		fmt.Println("error running program:", err.Error())
		os.Exit(1)
	}
}
