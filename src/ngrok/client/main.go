package client

import (
	"errors"
	"fmt"
	"github.com/inconshreveable/mousetrap"
	"math/rand"
	"ngrok/log"
	"ngrok/util"
	"os"
	"runtime"
	"time"
)

func init() {
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Don't double-click ngrok!")
			fmt.Println("You need to open cmd.exe and run it from the command line!")
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}
	}
}

func Main() {
	MainWithArgs(os.Args)
}

func parseCommandLine(arg0 string, command string) ([]string, error) {
	var args []string
	args = append(args, arg0)
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", command))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

func MainWithCommandLine(arg0 string, command string) int {
	var args, err = parseCommandLine(arg0, command)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	MainWithArgs(args)
	return 0
}

func MainWithCommandLineSync(arg0 string, command string) int {
	var args, err = parseCommandLine(arg0, command)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "ngrok exiting: %v\n", err)
		}
	}()
	MainWithArgs(args)
	return 0
}

func MainWithCommandLineAsync(arg0 string, command string) int {
	go func() {
		MainWithCommandLineSync(arg0, command)
	}()
	return 0
}

func MainWithArgs(args []string) {
	// parse options
	opts, err := ParseArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set up logging
	log.LogTo(opts.logto, opts.loglevel)

	// read configuration file
	config, err := LoadConfiguration(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		fmt.Printf("Couldn't securely seed the random number generator!")
		os.Exit(1)
	}
	rand.Seed(seed)

	NewController().Run(config)
}

func RunWithConfigFromFile(path string) {
	var opts = new(Options)
	opts.config = path
	opts.command = "start-all"
	// read configuration file
	config, err := LoadConfiguration(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		fmt.Printf("Couldn't securely seed the random number generator!")
		os.Exit(1)
	}
	rand.Seed(seed)

	NewController().Run(config)
}
