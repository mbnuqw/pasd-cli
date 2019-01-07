package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	con "github.com/mbnuqw/con-go"
)

var passAlph = [55]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'x',
	'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L',
	'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'2', '3', '4', '5', '6', '7', '8', '9', 'Y', 'Z', 'y',
}

// Command gen
func onGen(c *con.Client, args []string) {
	if len(args) < 1 {
		fmt.Printf(wrongInMsg() + usageMsg() + commandGenMsg())
		return
	}

	switch {
	case args[0] == "pass":
		genPass(c, args[1:])
	default:
		fmt.Printf(wrongInMsg() + usageMsg() + commandGenMsg())
		return
	}
}

// Generate password and put it in clipboard
func genPass(c *con.Client, args []string) {
	// Get desired pass-len
	var passLen int
	var err error
	if len(args) > 0 {
		passLen, err = strconv.Atoi(args[len(args)-1])
		if err != nil || passLen < 1 {
			passLen = 12
		}
	} else {
		passLen = 12
	}

	// Generate password
	var password bytes.Buffer
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < passLen; i++ {
		rd := rnd.Uint32()
		password.WriteRune(passAlph[rd%55])
	}

	// Put value in clipboard
	err = clipboard.WriteAll(password.String())
	if err != nil {
		PrintError("Cannot put value in clipboard")
		return
	}
}
