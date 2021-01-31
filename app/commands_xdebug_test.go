package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asannikov/jumper/app/command"
	"github.com/urfave/cli/v2"
)

func TestXdebugFpmEnableCase1(t *testing.T) {
	cliApp, jcfg, opt := jumperMainAppTest()
	opt.execCommand = func(eo command.ExecOptions, c *cli.App) error {
		fmt.Printf("\ncommand: %s\n\n", eo.GetCommand()+" "+strings.Join(eo.GetArgs(), " "))
		return nil
	}
	cliApp.Commands = commandList(jcfg.cfg, jcfg.dlg, opt)
	cliApp.Run([]string{"jumper", "xdebug:fpm:enable"})
}
