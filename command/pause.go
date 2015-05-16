package command

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/errwrap"

	"github.com/advanderveer/timer/model"
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
	return fmt.Sprintf("<description>")
}

func (c *Pause) Usage() string {
	return "<usage>"
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

	m := model.New(dir)
	info, err := m.ReadDaemonInfo()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Failed to get Daemon address: {{err}}"), err)
	}

	client := NewClient(info)
	err = client.Call("timer.pause")
	if err != nil {
		if err == ErrDaemonDown {
			return errwrap.Wrapf(fmt.Sprintf("No timer appears to be running for '%s': {{err}}", dir), err)
		} else {
			return err
		}
	}

	fmt.Println("Timer paused")
	return nil
}
