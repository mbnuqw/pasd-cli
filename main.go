package main

import (
	"fmt"
	"os"
	"time"

	con "github.com/mbnuqw/con-go"
)

func main() {
	// Setup config
	conf := InitConfig()

	// Connect to core service
	client := con.Client{}
	err := client.Connect(conf.CoreAddress, "pasd-cli")
	if err != nil {
		PrintError(err.Error())
		return
	}
	defer client.Disconnect()

	// Handle commands
	args := os.Args[1:]
	switch {
	case len(args) == 0:
		onHelp()
	case len(args) == 1 && args[0] == "help":
		onHelp()

	case len(args) == 1 && args[0] == "keys":
		listKeys(&client, args[1:])
	case len(args) == 1 && args[0] == "secrets":
		listSecrets(&client, args[1:])

	case len(args) > 1 && args[0] == "key":
		addKey(&client, args[1:])
	case len(args) > 1 && args[0] == "secret":
		addSecret(&client, args[1:])

	case len(args) > 1 && args[0] == "remove":
		onRemove(&client, args[1:])

	case args[0] == "gen":
		onGen(&client, args[1:])

	default:
		onGet(&client, args[:])
	}

	// Wait before exit
	time.Sleep(100 * time.Millisecond)
}

// Print general help message
func onHelp() {
	fmt.Printf(descMsg() + commandsMsg())
}

// Write out help message and exit.
func wrongCommand() {
	fmt.Printf(wrongInMsg() + usageMsg())
	os.Exit(1)
}
