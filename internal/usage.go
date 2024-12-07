package internal

import "fmt"

func CustomUsage() {
	fmt.Printf(`Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help
	
Options:
  --help       Show this screen.
  --port N     Port number.`)
	fmt.Print("\n")
}
