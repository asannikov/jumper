package app

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestXDebug(t *testing.T) {
	os.Args = []string{"jumper", "xdebug:fpm:enable"}

	c := &cli.App{}

	jcfg := &jumperAppTest{}
	JumperAppTest(c, jcfg)

	c.Run(os.Args)
}
