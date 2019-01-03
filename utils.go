package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var alph = [64]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_',
}

// UID generates simple uid string
func UID() string {
	// Get time part
	ns := uint64(time.Now().UnixNano())

	// Get rand part
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rd := rnd.Uint64()

	// Result string
	var output bytes.Buffer

	// Rand part
	for i := 0; i < 7; i++ {
		output.WriteRune(alph[rd&63])
		rd >>= 6
	}

	// Time part
	for i := 0; i < 5; i++ {
		output.WriteRune(alph[ns&63])
		ns >>= 6
	}

	return output.String()
}

// Remove last rune from string
func backspace(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return string(r[:len(r)-1])
}

// Ask user to input passwords
func askPV(keys []listedKey) (map[string]string, error) {
	passValues := map[string]string{}

	// Filter "text" keys
	passwords := []listedKey{}
	for i := range keys {
		if keys[i].Type == "text" {
			passwords = append(passwords, keys[i])
		}
	}

	if len(passwords) == 0 {
		return passValues, nil
	}

	passReader := bufio.NewReader(os.Stdin)
	passwordIndex := 0

	// disable input buffering, set min chars for complete read
	// and do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	// Hint
	fmt.Println(WithColors("|_>Enter passwords:|x|"))

	for {
		passName := passwords[passwordIndex].Name

		fmt.Printf(WithColors("|g>❯❯❯ |b>%s|x| "), passName)
		value := ""
		for {
			input, _, err := passReader.ReadRune()
			if err != nil {
				return nil, err
			}

			if input == '\n' {
				fmt.Print("\n")
				break
			}

			if input == 127 {
				value = backspace(value)
			} else {
				value += string(input)
			}

			fmt.Print("\r")
			fmt.Printf(getInputPrefix(len(value))+WithColors(" |b>%s|x| "), passName)
		}

		passValues[passName] = strings.TrimSpace(value)

		passwordIndex++
		if passwordIndex == len(passwords) {
			break
		}
	}

	return passValues, nil
}

// Ask user to input passwords name and values
func askPNV() (map[string]string, error) {
	passValues := map[string]string{}
	passReader := bufio.NewReader(os.Stdin)
	buffered := 0

	// disable input buffering, set min chars for complete read
	// and do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	// Hint
	fmt.Println(WithColors("|_>Enter passwords (Name Value):|x|"))

	for {
		fmt.Printf(WithColors("|g>❯❯❯ |x|"))
		name := ""
		space := false
		pass := ""
		for {
			input, _, err := passReader.ReadRune()
			if err != nil {
				return nil, err
			}

			// Check if there are any other runes
			// in buffer after readed one, which can be
			// possible if typing unprintable chars (arrows)
			buffered = passReader.Buffered()
			if buffered > 0 {
				_, err := passReader.Discard(buffered)
				if err != nil {
					return nil, err
				}
			}

			// Skip unsupported keys
			if input == 27 {
				continue
			}

			if input == '\n' {
				if name == "" || pass == "" {
					fmt.Print("\r                                    \n")
				} else {
					fmt.Print("\r                                    \r")
				}
				break
			}

			if !space {
				if input == 32 {
					space = true
				} else if input == 127 {
					fmt.Print("\r    " + strings.Repeat(" ", len([]rune(name))))
					name = backspace(name)
				} else {
					name += string(input)
				}
			} else {
				if input == 127 {
					if pass == "" {
						space = false
					}
					pass = backspace(pass)
				} else {
					pass += string(input)
				}
			}

			fmt.Print("\r")
			if space {
				fmt.Printf(getInputPrefix(len(pass))+WithColors(" |b>%s|x| "), name)
			} else {
				fmt.Printf(getInputPrefix(len(pass))+WithColors(" |b>%s|x|"), name)
			}
		}

		if string(pass) == "" {
			break
		}

		passValues[name] = pass
	}

	return passValues, nil
}

// Check error
func check(err error) bool {
	if err != nil {
		PrintError(err.Error())
		return true
	}
	return false
}
