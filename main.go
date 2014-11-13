package main

import (
	"fmt"
	"os"

	nanoAPI "github.com/nanobox-core/api-client-go"
	// "github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

const Version = "0.0.1"

type (

	// CLI represents the Nanobox CLI. It has a version, a Nanobox API client
	// and a map of all the commands it responds to
	CLI struct {
		version   string
		apiClient *nanoAPI.Client
		commands  map[string]Command
	}
)

// main creates a new CLI and then checks to see if authentication is needed. If
// no authentication is required it will attempt to run the provided command
func main() {

	//
	cli := &CLI{
		version:   Version,
		apiClient: nanoAPI.NewClient(),
		commands:  Commands,
	}

	// cli.apiClient.APIURL     = "localhost:8080"
	// cli.apiClient.APIVersion = "api"
	// cli.apiClient.AuthToken = ""

	// run the CLI
	cli.run()

}

// run attempts to run a CLI command. If no flags are passed (only the program
// is run) it will default to printing the CLI help text. It takes a help flag
// for printing the CLI help text. It takes a version flag for displaying the
// current version. It takes an app flag to indicate which app to run the command
// on (otherwise it wll attempt to find an app associated with the current directory).
// It also takes a debug flag (which must be passed last), that will display all
// request/response output for any API call the CLI makes.
func (cli *CLI) run() {

	// command line args w/o program
	args := os.Args[1:]

	// if only program is run, print help by default
	if len(args) <= 0 {
		cli.Help()

		// parse command line args
	} else {

		// it's safe to assume that args[0] is the command we want to run, or one of
		// our 'shortcut' flags that we'll catch before trying to run the command.
		command := args[0]

		// check for 'global' commands
		switch command {

		// Check for help shortcuts
		case "-h", "--help", "help":
			cli.Help()

		// Check for version shortcuts
		case "-v", "--version", "version":
			ui.CPrintln("[yellow]Version " + cli.Version() + "[reset]")

		// we didn't find a 'shortcut' flag, so we'll continue parsing the remaining
		// args looking for a command to run.
		default:

			// if we find a valid command we run it
			if val, ok := cli.commands[command]; ok {

				// args[1:] will be our remaining subcommand or flags after the intial command.
				// This value could also be 0 if running an alias command.
				opts := args[1:]

				// assume they wont be passing an app
				fApp := ""

				//
				if len(opts) >= 1 {
					switch opts[0] {

					// Check for help shortcuts
					case "-h", "--help", "help":
						cli.commands[command].Help()
						os.Exit(0)

					// Check for app flag, set fApp and strip out the flag and app
					case "-a", "--app":
						fApp = opts[1]
						opts = opts[2:]
					}
				}

				// before we run the command we'll check to see if debug mode needs to
				// be enabled. If so, enable it and strip off the flag.
				if args[len(args)-1] == "--debug" {
					cli.apiClient.Debug = true

					opts = opts[:len(opts)-1]
				}

				// run the command
				val.Run(fApp, opts, cli.apiClient)

				// no valid command found
			} else {
				fmt.Printf("'%v' is not a valid command. Type 'pagoda' for available commands\n and usage.", command)
				os.Exit(1)
			}
		}
	}
}
