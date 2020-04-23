package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	app := cli.App{
		Name: "mongobackup",
		Usage: "uses mongodump to create a backup from a mongo db",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "uri",
				Usage:    "The db connection string to create a backup for.",
				EnvVars:  []string{"MONGO_BACKUP_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dir",
				Usage:    "the path to the folder to store the backup",
				EnvVars:  []string{"MONGO_BACKUP_DIR"},
				Required: true,
			},
		},

		Action: func(ctx *cli.Context) error {
			uri := ctx.String("uri")
			dir := ctx.String("dir")

			if s, err := os.Stat(dir); err != nil || s.IsDir() == false {
				return errors.New(fmt.Sprintf("%v is not a dir", dir))
			}
			path, err := triggerDump(uri, dir)
			if err == nil {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("Successfully dumped db to %s", path))
			}
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}


func triggerDump(uri, dir string) (string, error) {
	if mongodumpExists() == false {
		return "", errors.New("<mongodump> not found but required")
	}
	path, err := createDir(dir)
	if err != nil {
		return "", err
	}

	args := []string{
		fmt.Sprintf("--out=\"%v\"", path),
		fmt.Sprintf("--uri=\"%v\"", uri),
	}

	cmd := exec.Command("mongodump", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		wipeDir(path)
		return "", err
	}
	return path, nil
}

func mongodumpExists() bool {
	_,err := exec.LookPath("mongodump")

	return err == nil
}

func createDir(out string) (string, error) {
	dir := fmt.Sprintf("%v%vmongo_backup_%v", out, string(os.PathSeparator), time.Now().Unix())

	return dir, os.Mkdir(dir, 0755)
}

func wipeDir(path string) {
	os.RemoveAll(path)
}