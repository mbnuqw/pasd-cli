package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	con "github.com/mbnuqw/con-go"
)

type addKeyArgs struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Group     string            `json:"group,omitempty"`
	Value     string            `json:"value"`
	Passwords map[string]string `json:"passwords"`
}

type addKeyAns struct {
	Error string `json:"error,omitempty"`
}

type addSecretArgs struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Value     string            `json:"value"`
	Login     string            `json:"login,omitempty"`
	URL       string            `json:"url,omitempty"`
	Passwords map[string]string `json:"passwords"`
}

type addSecretAns struct {
	Error string `json:"error,omitempty"`
}

// addKey - requests existed keys, parses input
// and asks user for passwords to auth and reencrypt db
func addKey(c *con.Client, args []string) {
	// Check count of args and types
	if len(args) < 1 || len(args) > 2 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandAddMsg())
		return
	}

	// Get keys
	keys, err := getKeys(c)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Parse new key's name [and group]
	keyNameGroupStr := args[0]
	keyNameGroup := strings.Split(keyNameGroupStr, ":")
	keyName := keyNameGroup[0]
	var keyGroup string
	if len(keyNameGroup) == 2 {
		keyGroup = keyNameGroup[1]
	}

	// Parse last argument
	var keyValue string
	var keyType string
	if len(args) == 2 && strings.HasPrefix(args[1], "/") {
		keyType = "file"
		keyValue = args[1]
	} else {
		keyType = "text"
		newKeyPlaceholder := listedKey{Name: keyName, Type: "text"}
		keys = append([]listedKey{newKeyPlaceholder}, keys...)
	}

	// Get passwords
	passValues, err := askPV(keys)
	if err != nil {
		PrintError(err.Error())
		return
	}
	if keyType == "text" {
		keyValue = passValues[keyName]
	}

	// Serialize arguments
	argsJSON, err := json.Marshal(addKeyArgs{
		Type:      keyType,
		Name:      keyName,
		Group:     keyGroup,
		Value:     keyValue,
		Passwords: passValues,
	})
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandAddMsg())
		return
	}

	// Make request
	ansChan := make(chan con.Msg, 1)
	go c.Req("add-key", argsJSON, ansChan)

	// Wait
	fmt.Print(WithColors("|b>Loading...|x|"))
	ansJSON := <-ansChan
	fmt.Print("\r")

	// Parse answer
	var ans addKeyAns
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
	fmt.Println(WithColors("|g>✔ Key added!|x|"))
}

// addSecret - request extisted keys, check arguments,
// ask user to input passwords and finally send req
// to core.
func addSecret(c *con.Client, args []string) {
	// Check count of args and types
	if len(args) < 1 || len(args) > 3 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandAddMsg())
		return
	}

	// Get keys
	keys, err := getKeys(c)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Parse secret info
	secretName := args[0]
	var secretLogin string
	var secretURL string
	if len(args) == 2 {
		secretURL = args[1]
	}
	if len(args) == 3 {
		secretURL = args[1]
		secretLogin = args[2]
	}

	// Get secret value
	var secretValue string
	fmt.Print(WithColors("|g>❯❯❯ |x|"))
	passReader := bufio.NewReader(os.Stdin)
	input, _, err := passReader.ReadLine()
	if check(err) {
		return
	}
	secretValue = string(input)

	// Cleanup
	cleanLen := len([]rune(secretValue))
	fmt.Print("\033[1A")
	fmt.Print("\r    " + strings.Repeat(" ", cleanLen) + "\r")

	// Check type of secret
	var secretType string
	valueIsAbsPath := strings.HasPrefix(secretValue, "/")
	valueIsRelPath := strings.HasPrefix(secretValue, "./")
	if valueIsAbsPath || valueIsRelPath {
		secretType = "file"
		if valueIsRelPath {
			secretValue, err = filepath.Abs(secretValue)
			if err != nil {
				PrintError(err.Error())
			}
		}
	} else {
		secretType = "text"
	}

	// Get passwords
	passValues, err := askPV(keys)
	if err != nil {
		PrintError(err.Error())
		return
	}

	// Serialize arguments
	argsJSON, err := json.Marshal(addSecretArgs{
		Type:      secretType,
		Name:      secretName,
		Value:     secretValue,
		URL:       secretURL,
		Login:     secretLogin,
		Passwords: passValues,
	})
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandAddMsg())
		return
	}

	// Send
	ansChan := make(chan con.Msg, 1)
	go c.Req("add-secret", argsJSON, ansChan)

	// Wait
	fmt.Print(WithColors("|b>Loading...|x|"))
	ansJSON := <-ansChan
	fmt.Print("\r")

	// Parse answer
	var ans addSecretAns
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
	fmt.Println(WithColors("|g>✔ Secret added!|x|"))
}
