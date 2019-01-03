package main

import (
	"encoding/json"
	"fmt"
	"time"

	con "github.com/mbnuqw/con-go"
)

type listedKey struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Date  int64  `json:"date"`
}

type listKeysAns struct {
	Keys  []listedKey `json:"keys"`
	Error string      `json:"error,omitempty"`
}

type listedSecret struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Login string `json:"login"`
	Date  int64  `json:"date"`
}

type listSecretsAns struct {
	Secrets []listedSecret `json:"secrets"`
	Error   string         `json:"error,omitempty"`
}

func getKeys(c *con.Client) ([]listedKey, error) {
	ansChan := make(chan con.Msg, 1)
	go c.Req("list-keys", nil, ansChan)
	ansMsg := <-ansChan

	var ans listKeysAns
	err := json.Unmarshal(ansMsg.Body, &ans)
	if err != nil {
		return nil, err
	}

	if ans.Error != "" {
		return nil, NewCoreError(ans.Error)
	}

	return ans.Keys, nil
}

// List keys
func listKeys(c *con.Client, args []string) {
	ansChan := make(chan con.Msg, 1)
	go c.Req("list-keys", nil, ansChan)
	ansMsg := <-ansChan
	var ans listKeysAns
	err := json.Unmarshal(ansMsg.Body, &ans)
	if err != nil {
		PrintError(err.Error())
		return
	}

	if len(ans.Keys) == 0 {
		fmt.Println("")
		fmt.Println(WithColors("|_> Nothing...|x|"))
		return
	}

	// Get groups
	groups := []string{}
Groups:
	for i := range ans.Keys {
		key := ans.Keys[i]
		for g := range groups {
			if groups[g] == key.Group {
				continue Groups
			}
		}

		groups = append(groups, key.Group)
	}

	// Print out
	fmt.Println(" ")
	fmt.Printf(WithColors("|B> |w>%-20s  %-8s  %-16s|x|\n"), "Group/Name", "Type", "Date")
	for g := range groups {
		group := groups[g]
		fmt.Printf(WithColors("|_> %s:|x|\n"), group)
		for i := range ans.Keys {
			k := ans.Keys[i]
			date := time.Unix(k.Date, 0).Format("2006.01.02 15:04")
			if k.Group != group {
				continue
			}
			fmt.Printf(WithColors("     |g>%-16s  %-8s  %-16s|x|\n"), k.Name, k.Type, date)
		}
	}
}

// List secrets
func listSecrets(c *con.Client, args []string) {
	ansChan := make(chan con.Msg, 1)
	go c.Req("list-secrets", nil, ansChan)
	ansMsg := <-ansChan
	var ans listSecretsAns
	err := json.Unmarshal(ansMsg.Body, &ans)
	if err != nil {
		PrintError(err.Error())
		return
	}

	if len(ans.Secrets) == 0 {
		fmt.Println("")
		fmt.Println(WithColors("|_> Nothing...|x|"))
		return
	}

	fmt.Println("")
	fmt.Printf(WithColors("|B>|w> %-16s  %-13s  %-20s  %-19s|x|\n"), "Date", "Name", "URL", "Login")
	for i := range ans.Secrets {
		s := ans.Secrets[i]
		date := time.Unix(s.Date, 0).Format("2006.01.02 15:04")
		fmt.Printf(WithColors("|_> %-16s  |g>%-13s  %-20s  %-19s|x|\n"), date, s.Name, s.URL, s.Login)
	}
}
