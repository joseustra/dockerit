package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
)

// A Config is the configuration for the container
type Config struct {
	Container Container
}

// A Container is the configuration for the app
type Container struct {
	Name   string
	Port   string
	Link   string
	Volume string
	Image  string
}

func loadConfig() (container *Container, pwd string) {
	pwd, _ = os.Getwd()
	dir := strings.Split(pwd, "/")
	appName := dir[len(dir)-1]

	data, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.dockerit/" + appName + ".yml")
	config := new(Config)

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	container = &config.Container

	return container, pwd
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

	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		fmt.Printf("docker build out | %s\n", stderr.String())
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("docker build out | %s\n", stderr.String())
		return err
	}
	return err
}

func cleanContainer(name string) {
	cmdName := "docker"

	cmdArgs := []string{"ps", "-a", "-f", "name=" + name}
	cmd := exec.Command(cmdName, cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error to exec Command", err)
		os.Exit(1)
	}

	r := regexp.MustCompile(name)
	if r.Match([]byte(string(out))) {
		cmdArgs = []string{"rm", "-f", name}
		cmd = exec.Command(cmdName, cmdArgs...)
		err = execCommand(cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to exec Command", err)
			os.Exit(1)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dockerit"
	app.Usage = "Run Docker commands easily"

	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build an image",
			Action: func(c *cli.Context) {
				container, _ := loadConfig()

				cmdName := "docker"
				cmdArgs := []string{"build", "-t", container.Image, "."}

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
				container, pwd := loadConfig()
				cleanContainer(container.Name)

				cmdName := "docker"
				cmdArgs := []string{"run", "--name", container.Name, "-a", "stdout", "-a", "stderr", "-p", container.Port, "--link", container.Link, "-v", pwd + ":/go/src/app", container.Image}

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
				container, pwd := loadConfig()
				cleanContainer(container.Name)

				cmdName := "docker"
				cmdArgs := []string{"run", "--name", container.Name, "-a", "stdout", "-a", "stderr", "-p", container.Port, "--link", container.Link, "-v", pwd + ":/go/src/app", container.Image}

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
