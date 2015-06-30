package main

import (
	"fmt"
	"os"

	"github.com/timeglass/glass/_vendor/github.com/codegangsta/cli"

	"github.com/timeglass/glass/command"
)

type Command interface {
	Name() string
	Description() string
	Usage() string
	Run(c *cli.Context) error
	Action() func(ctx *cli.Context)
	Flags() []cli.Flag
}

var Version = "0.0.0"
var Build = "gobuild"

func main() {
	app := cli.NewApp()
	app.Author = "Ad van der Veer"
	app.Email = "advanderveer@gmail.com"
	app.Name = "Timeglass"
	app.Usage = "Automated time tracking for code repositories"
	app.Version = fmt.Sprintf("%s (%s)", Version, Build)

	cmds := []Command{
		command.NewInstall(), //install daemon as a service on various platforms
		command.NewInit(),    //write hooks, create timer and pull time data
		command.NewStart(),   //create timer for current directory
		// command.NewPause(),
		// command.NewStatus(),
		command.NewStop(), //remove timer for current directory
		command.NewPush(),
		command.NewPull(),
		// command.NewLap(),
		command.NewPunch(),
		command.NewSum(),
	}

	for _, c := range cmds {
		app.Commands = append(app.Commands, cli.Command{
			Name:        c.Name(),
			Usage:       c.Usage(),
			Action:      c.Action(),
			Description: c.Description(),
			Flags:       c.Flags(),
		})
	}

	app.Run(os.Args)
}
