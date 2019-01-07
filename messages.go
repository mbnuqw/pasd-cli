package main

// DescMsg does nothing (for now)
func descMsg() string {
	return WithColors(`
  |w>This is CLI util for working with 'pasd' - secrets management tool.|x|
`)
}

func commandsMsg() string {
	return WithColors(`
  |y>Commands:|c>
    |c>help|_> - print this message
`) + commandListMsg() + commandAddMsg() + commandRemoveMsg() + commandGetMsg() + commandGenMsg()
}

func commandListMsg() string {
	return WithColors(`
    |c>keys|_> - list all keys
    |c>secrets|_> - list all secrets
`)
}

func commandAddMsg() string {
	return WithColors(`
    |c>key |g><Name>[:<Group>] [Path]|_> - add new key
      |_>pasd key Master - add password key
      |_>pasd key Some /path/to/key - add file key
      |_>pasd key A:Common ┬─ add grouped keys
      |_>pasd key B:Common ┘

    |c>secret|g> <Name> [URL] [Login]|_> - add new secret
      |_>pasd secret Example example.com login
      |_>  >>> 123456
      |_>pasd secret SomeFile
      |_>  >>> /path/to/file
`)
}

func commandRemoveMsg() string {
	return WithColors(`
    |c>remove key |g><Name>|_> - remove key
    |c>remove secret |g><Name>|_> - remove secret
`)
}

func commandGetMsg() string {
	return WithColors(`
    |g><Search query>|_> - get some secret
      |_>pasd login gmail
      |_>pasd goo / Master 123 AnotherKey 456
`)
}

func commandGenMsg() string {
	return WithColors(`
    |c>gen pass|g> <Length>|_> - generate password with given length
`)
}

// WrongInMsg prints error about wrong cli params
func wrongInMsg() string {
	return WithColors(`|r>  Wrong input|x|

`)
}

func usageMsg() string {
	return WithColors(`  |y>Usage:`)
}

func getInputPrefix(i int) string {
	switch i % 5 {
	case 0:
		return WithColors("|g>❯❯❯|x|")
	case 1:
		return WithColors("|g> ❯❯|x|")
	case 2:
		return WithColors("|g>  ❯|x|")
	case 3:
		return WithColors("|g>❯  |x|")
	default:
		return WithColors("|g>❯❯ |x|")
	}
}
