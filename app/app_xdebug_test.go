package app

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestXDebug(t *testing.T) {
	os.Args = []string{"xdebug:fpm:enable"}

	c := &cli.App{}

	JumperAppTest(c)

	c.Run(os.Args)
}
