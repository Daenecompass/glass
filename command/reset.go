package command

import (
	"fmt"
	"os"

	"github.com/timeglass/glass/_vendor/github.com/codegangsta/cli"
	"github.com/timeglass/glass/_vendor/github.com/hashicorp/errwrap"
)

type Reset struct {
	*command
}

func NewReset() *Reset {
	return &Reset{newCommand()}
}

func (c *Reset) Name() string {
	return "reset"
}

func (c *Reset) Description() string {
	return fmt.Sprintf("...")
}

func (c *Reset) Usage() string {
	return "..."
}

func (c *Reset) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *Reset) Action() func(ctx *cli.Context) {
	return c.command.Action(c.Run)
}

func (c *Reset) Run(ctx *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return errwrap.Wrapf("Failed to fetch current working dir: {{err}}", err)
	}

	c.Printf("Resetting timer to 0s...")

	client := NewClient()
	err = client.ResetTimer(dir)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to reset timer: {{err}}"), err)
	}

	c.Printf("Timer is reset!")
	return nil
}
