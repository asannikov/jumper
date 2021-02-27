package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

type magentoGlobalConfig interface {
	SaveContainerNameToProjectConfig(string) error
	GetProjectMainContainer() string
	SaveDockerProjectPath(string) error
	GetProjectDockerPath() string
}

type magentoDialog interface {
	SetMainContaner([]string) (int, string, error)
	DockerProjectPath(string) (string, error)
}

type magentoOptions interface {
	GetExecCommand() func(ExecOptions, *cli.App) error
	GetCommandLocation() func(string, string) (string, error)
	GetInitFunction() func(bool) string
	GetContainerList() ([]string, error)
	CheckMagentoBin(containerName string, magentoBin string) (bool, error)
}

func callMagentoCommanBin(cfg magentoGlobalConfig, d magentoDialog, options magentoOptions) *cli.Command {
	execCommand := options.GetExecCommand()
	initf := options.GetInitFunction()

	return &cli.Command{
		Name:    "bin/magento",
		Aliases: []string{"bm"},
		Usage:   "It looks for bin/magento command and call it",
		Action: func(c *cli.Context) error {
			initf(true)

			var err error
			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			if err = defineProjectDockerPath(cfg, d, "/var/www/html/"); err != nil {
				return err
			}

			paths := []string{
				"bin/magento",
				"html/bin/magento",
				"source/bin/magento",
				"src/bin/magento",
			}

			var magentoBinSource string
			var status bool
			for _, path := range paths {
				p := strings.TrimRight(cfg.GetProjectDockerPath(), string(os.PathSeparator)) + string(os.PathSeparator) + path

				if status, err = options.CheckMagentoBin(cfg.GetProjectMainContainer(), p); err != nil {
					return err
				}

				if status {
					magentoBinSource = p
					break
				}
			}

			if magentoBinSource == "" {
				return fmt.Errorf("Cannot find magento root folder. Searched for: %s", paths)
			}

			var args []string

			args = append(args, []string{"exec", "-it", cfg.GetProjectMainContainer(), magentoBinSource}...)
			args = append(args, c.Args().Slice()...)

			eo := ExecOptions{
				command: "docker",
				args:    args,
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}
}

func callMagentoCommandMageRun(cfg magentoGlobalConfig, d magentoDialog, options magentoOptions) *cli.Command {
	execCommand := options.GetExecCommand()
	commandLocation := options.GetCommandLocation()
	initf := options.GetInitFunction()

	return &cli.Command{
		Name:    "magerun",
		Aliases: []string{"mr"},
		Usage:   "It looks for magerun2 command and call it",
		Action: func(c *cli.Context) error {
			initf(true)

			var err error
			var cl []string

			if cl, err = options.GetContainerList(); err != nil {
				return err
			}

			if err = defineProjectMainContainer(cfg, d, cl); err != nil {
				return err
			}

			var mrPath string

			if mrPath, err = commandLocation(cfg.GetProjectMainContainer(), "n98-magerun2.phar"); err != nil {
				return err
			}

			var args []string

			args = append(args, []string{"exec", "-it", cfg.GetProjectMainContainer(), mrPath}...)
			args = append(args, c.Args().Slice()...)

			eo := ExecOptions{
				command: "docker",
				args:    args,
				tty:     true,
				detach:  true,
			}

			return execCommand(eo, c.App)
		},
	}
}

// CallMagentoCommand runs copyright dialog
func CallMagentoCommand(cfg magentoGlobalConfig, d magentoDialog, options magentoOptions) *cli.Command {

	return &cli.Command{
		Name:    "magento",
		Aliases: []string{"m"},
		Usage:   "Call magento command bin/magento or magerun. This command has subcommands. Call jumper magento for more details.",
		Subcommands: []*cli.Command{
			callMagentoCommanBin(cfg, d, options),
			callMagentoCommandMageRun(cfg, d, options),
		},
	}
}

//CheckMagentoBin checks magento bin path
func CheckMagentoBin(containerName string, magentoBin string) (bool, error) {
	if len(containerName) == 0 {
		return false, errors.New("Container is not defined")
	}

	if len(magentoBin) == 0 {
		return false, errors.New("bin/magento is not defined")
	}

	var args = []string{"exec", "-i", containerName, "sh"}

	cmd := exec.Command("docker", args...)

	r, w := io.Pipe()

	cmd.Stdout = w
	cmd.Stderr = w

	var stdin io.WriteCloser
	var err error

	if stdin, err = cmd.StdinPipe(); err != nil {
		return false, fmt.Errorf("error when piping stdin: %w", err)
	}

	defer stdin.Close()

	scanMergedOut := bufio.NewScanner(r)

	fmt.Fprintln(stdin, `[ -f "`+magentoBin+`" ] && echo "yes" || echo "no"`)

	if err = cmd.Start(); err != nil {
		return false, fmt.Errorf("Cannot start command: %w", err)
	}

	if !scanMergedOut.Scan() {
		return false, errors.New("nothing found in Scanner output")
	}

	if scanMergedOut.Err() != nil {
		return false, fmt.Errorf("error while reading scanMergedOut: %s", scanMergedOut.Err())
	}

	return strings.Trim(string(scanMergedOut.Bytes()), " ") == "yes", nil
}
