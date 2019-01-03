package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/atotto/clipboard"
	con "github.com/mbnuqw/con-go"
)

type getSecretArgs struct {
	Query     []string          `json:"query"`
	Passwords map[string]string `json:"passwords"`
}

type getSecretAns struct {
	Secret []byte `json:"secret,omitempty"`
	Type   string `json:"type,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Command get
func onGet(c *con.Client, args []string) {
	if len(args) == 0 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGetMsg())
	}

	// Get password if needed
	passwords, err := askPNV()
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGetMsg())
		return
	}

	// Parse arguments for file to output
	var outputFilePath string
	if len(args) > 1 {
		theLastArg := args[len(args)-1]
		isAbsPath := strings.HasPrefix(theLastArg, "/")
		isRelPath := strings.HasPrefix(theLastArg, "./")
		if isAbsPath || isRelPath {
			outputFilePath = theLastArg
			args = args[:len(args)-1]
		}
	}

	// Prepare args for request
	argsJSON, err := json.Marshal(getSecretArgs{
		Query:     args,
		Passwords: passwords,
	})
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGetMsg())
		return
	}

	// Make request and wait
	ansChan := make(chan con.Msg, 1)
	go c.Req("get-secret", argsJSON, ansChan)

	// Wait
	fmt.Print(WithColors("|b>Loading...|x|"))
	ansJSON := <-ansChan
	fmt.Print("\r")

	// Handle answer
	var ans getSecretAns
	err = json.Unmarshal(ansJSON.Body, &ans)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Output
	if ans.Error != "" {
		PrintError(ans.Error)
		return
	}

	// Write secret to file
	if outputFilePath != "" {
		err = ioutil.WriteFile(outputFilePath, ans.Secret, 0600)
		fmt.Println(WithColors("|g>✔ Secret saved to file|x|"))
		return
	}

	// Put secret to clipboard
	if outputFilePath == "" {
		if ans.Type == "file" {
			PrintError("You need to provide output file path\n" +
				WithColors("$ pasd search query |B>PATH|N> |_>ԅ(≖‿≖ԅ)|x|"))
			return
		}
		err = clipboard.WriteAll(string(ans.Secret))
		fmt.Println(WithColors("|g>✔ Secret in clipboard|x|"))
	}
}
