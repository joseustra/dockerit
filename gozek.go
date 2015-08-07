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
	"github.com/codeskyblue/go-sh"
)

// A Config is the configuration for the app
type Config struct {
	Name    string
	Port    string
	Link    string
	Image   string
	Volume  string
	Command string
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

func pipeCommands(commands ...*exec.Cmd) []byte {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil
	}
	return final
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
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
					os.Exit(1)
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
					os.Exit(1)
				}

				err = cmd.Wait()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run an image",
			Action: func(c *cli.Context) {
				config, pwd := loadConfig()

				err := sh.Command("docker", "run", "--name", config.Name, "-a", "stdout", "-a", "stderr", "-i", "-t", "-p", config.Port, "--link",
					config.Link, "-v", pwd+":/go/src/"+config.Name, config.Image, config.Command).Run()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
