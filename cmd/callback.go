package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/hook"
)

const (
	// ErrExitCode represent error exit code
	ErrExitCode = 1

	// HookDir represent default hook directory
	HookDir = ".commitlint/hooks"
)

// VersionCallback is the version command callback
var VersionCallback = func() error {
	fmt.Println("commitlint")
	return nil
}

func initCallback(ctx *cli.Context) (retErr error) {
	// get user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create hooks dir
	hookDir := filepath.Join(homeDir, filepath.Clean(HookDir))
	err = os.MkdirAll(hookDir, os.ModePerm)
	if err != nil {
		return err
	}

	// create hook file
	err = hook.WriteToFile(hookDir)
	if err != nil {
		return err
	}

	isGlobal := ctx.Bool("global")
	return setGitConf(hookDir, isGlobal)
}

func lintCallback(ctx *cli.Context) error {
	confFilePath := ctx.String("config")
	fileInput := ctx.String("message")

	resStr, hasError, err := runLint(confFilePath, fileInput)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}

	if hasError {
		return cli.Exit(resStr, ErrExitCode)
	}

	fmt.Println(resStr)
	return nil
}

func hookCreateCallback(ctx *cli.Context) (retErr error) {
	err := hook.WriteToFile(".")
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}
	return nil
}

func configCreateCallback(ctx *cli.Context) error {
	isOnlyEnabled := ctx.Bool("enabled")
	err := config.DefaultConfToFile(isOnlyEnabled)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}
	return nil
}