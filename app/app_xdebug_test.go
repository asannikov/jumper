package app

import (
	"os"
	"testing"

	"github.com/asannikov/jumper/app/config"
	"github.com/asannikov/jumper/app/dialog"
	"github.com/urfave/cli/v2"
)

func TestXDebug(t *testing.T) {
	os.Args = []string{"jumper", "xdebug:fpm:enable"}

	c := &cli.App{}

	DLG := dialog.InitDialogFunctions()

	//DLG.SetMainContaner

	cfg := &config.Config{
		ProjectFile: confgFile,
	}

	jcfg := &jumperAppTest{
		dlg: &DLG,
		cfg: cfg,
		fs:  &FileSystem{},
	}

	JumperAppTest(c, jcfg)

	c.Run(os.Args)
}
