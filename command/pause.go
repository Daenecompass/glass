package command

import (
	"fmt"
	"os"

	"github.com/timeglass/glass/_vendor/github.com/codegangsta/cli"
	"github.com/timeglass/glass/_vendor/github.com/hashicorp/errwrap"
)

type Pause struct {
	*command
}

func NewPause() *Pause {
	return &Pause{newCommand()}
}

func (c *Pause) Name() string {
	return "pause"
}

func (c *Pause) Description() string {
	return fmt.Sprintf("...")
}

func (c *Pause) Usage() string {
	return "..."
}

func (c *Pause) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *Pause) Action() func(ctx *cli.Context) {
	return c.command.Action(c.Run)
}

func (c *Pause) Run(ctx *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return errwrap.Wrapf("Failed to fetch current working dir: {{err}}", err)
	}

	c.Printf("Pausing timer...")

	client := NewClient()
	err = client.PauseTimer(dir)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to pause timer: {{err}}"), err)
	}

	c.Printf("Done!")
	return nil
}
