package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
)

// A Config is the configuration for the app
type Config struct {
	Name   string
	Port   string
	Link   string
	Volume string
	Image  string
}

func loadConfig() (config *Config, pwd string) {
	pwd, _ = os.Getwd()
	dir := strings.Split(pwd, "/")
	appName := dir[len(dir)-1]

	data, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.gozek/" + appName + ".yml")
	config = new(Config)

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config, pwd
}

func execCommand(cmd *exec.Cmd) error {
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("docker build out | %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return err
	}
	return err
}

func main() {
	app := cli.NewApp()
	app.Name = "Gozek"
	app.Usage = "Run Docker commands easily"

	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build an image",
			Action: func(c *cli.Context) {
				config, _ := loadConfig()

				cmdName := "docker"
				cmdArgs := []string{"build", "-t", config.Image, "."}

				cmd := exec.Command(cmdName, cmdArgs...)
				err := execCommand(cmd)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error to exec Command", err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "up",
			Aliases: []string{"up"},
			Usage:   "up a container",
			Action: func(c *cli.Context) {
				config, pwd := loadConfig()

				cmdName := "docker"
				cmdArgs := []string{"run", "--name", config.Name, "-a", "stdout", "-a", "stderr", "-it", "-p", config.Port, "--link",
					config.Link, "-v", pwd + ":/go/src/app", config.Image}

				cmd := exec.Command(cmdName, cmdArgs...)
				err := execCommand(cmd)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error to exec Command", err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "run",
			Aliases: []string{"run"},
			Usage:   "run an command on the container",
			Action: func(c *cli.Context) {
				config, pwd := loadConfig()

				cmdName := "docker"
				cmdArgs := []string{"run", "--name", config.Name, "-a", "stdout", "-a", "stderr", "-i", "-t", "-p", config.Port, "--link",
					config.Link, "-v", pwd + ":/go/src/app", config.Image}

				for _, arg := range c.Args() {
					cmdArgs = append(cmdArgs, arg)
				}

				cmd := exec.Command(cmdName, cmdArgs...)
				err := execCommand(cmd)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error to exec Command", err)
					os.Exit(1)
				}
			},
		},
	}

	app.Run(os.Args)
}
