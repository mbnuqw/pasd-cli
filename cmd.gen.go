package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	con "github.com/mbnuqw/con-go"
)

type genSecretArgs struct {
	Name  string            `json:"name"`
	Login string            `json:"login,omitempty"`
	URL   string            `json:"url,omitempty"`
	Len   int               `json:"len,omitempty"`
	Keys  map[string]string `json:"keys,omitempty"`
}

type genSecretAns struct {
	Secret []byte `json:"secret,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Command gen
func onGen(c *con.Client, args []string) {
	if len(args) < 2 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGenMsg())
		return
	}

	secretName := args[0]
	var secretURL string
	var secretLogin string
	var secretLen int
	var err error
	if len(args) == 2 {
		secretLen, err = strconv.Atoi(args[1])
	}
	if len(args) == 3 {
		secretURL = args[1]
		secretLen, err = strconv.Atoi(args[2])
	}
	if len(args) == 4 {
		secretURL = args[1]
		secretLogin = args[2]
		secretLen, err = strconv.Atoi(args[3])
	}
	if err != nil {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGenMsg())
		return
	}
	// TODO text keys...
	textKeys := map[string]string{}

	argsJSON, err := json.Marshal(genSecretArgs{
		Name:  secretName,
		URL:   secretURL,
		Login: secretLogin,
		Len:   secretLen,
		Keys:  textKeys,
	})

	ansChan := make(chan con.Msg, 1)
	go c.Req("gen-secret", argsJSON, ansChan)
	ansJSON := <-ansChan
	var ans genSecretAns
	err = json.Unmarshal(ansJSON.Body, &ans)
	if err != nil {
		fmt.Println(" → Something goes wrong")
		return
	}
	fmt.Println(" → gen", ans)
}
