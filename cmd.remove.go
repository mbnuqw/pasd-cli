package main

import (
	"encoding/json"
	"fmt"

	con "github.com/mbnuqw/con-go"
)

type removeSecretArgs struct {
	Query     []string          `json:"query"`
	Passwords map[string]string `json:"passwords"`
}

type removeKeyArgs struct {
	Name      string            `json:"name"`
	Passwords map[string]string `json:"passwords"`
}

type removeKeyAns struct {
	Error string `json:"error,omitempty"`
}

type removeSecretAns struct {
	Error string `json:"error,omitempty"`
}

// Command remove
func onRemove(c *con.Client, args []string) {
	switch {
	case len(args) > 0 && args[0] == "key":
		removeKey(c, args[1:])
	case len(args) > 0 && args[0] == "secret":
		removeSecret(c, args[1:])
	default:
		fmt.Printf(wrongInMsg() + usageMsg() + commandRemoveMsg())
	}
}

// ...
func removeKey(c *con.Client, args []string) {
	if len(args) == 0 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandRemoveMsg())
		return
	}

	// Get keys
	passwords, err := getKeys(c)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Get passwords
	passValues, err := askPV(passwords)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Prepare arguments
	argsJSON, err := json.Marshal(removeKeyArgs{
		Name:      args[0],
		Passwords: passValues,
	})
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandRemoveMsg())
		return
	}

	// Make request
	ansChan := make(chan con.Msg, 1)
	go c.Req("remove-key", argsJSON, ansChan)

	// Wait
	fmt.Print(WithColors("|b>Loading...|x|"))
	ansJSON := <-ansChan
	fmt.Print("\r")

	// Parse answer
	var ans removeKeyAns
	err = json.Unmarshal(ansJSON.Body, &ans)
	if err != nil {
		fmt.Println(" → Something goes wrong")
		return
	}

	// Print result
	if ans.Error != "" {
		PrintError(ans.Error)
		return
	}
	fmt.Println(WithColors("|g>✔ Key successfully removed!|x|"))
}

func removeSecret(c *con.Client, args []string) {
	if len(args) == 0 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandRemoveMsg())
		return
	}

	// Get keys
	passwords, err := getKeys(c)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Get passwords
	passValues, err := askPV(passwords)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Prepare arguments
	argsJSON, err := json.Marshal(removeSecretArgs{
		Query:     args,
		Passwords: passValues,
	})
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandRemoveMsg())
		return
	}

	// Make request
	ansChan := make(chan con.Msg, 1)
	go c.Req("remove-secret", argsJSON, ansChan)

	// Wait
	fmt.Print(WithColors("|b>Loading...|x|"))
	ansJSON := <-ansChan
	fmt.Print("\r")

	// Parse answer
	var ans removeSecretAns
	err = json.Unmarshal(ansJSON.Body, &ans)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Print result
	if ans.Error != "" {
		PrintError(ans.Error)
		return
	}
	fmt.Println(WithColors("|g>✔ Secret successfully removed!|x|"))
}
