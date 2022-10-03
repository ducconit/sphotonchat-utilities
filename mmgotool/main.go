package main

import (
	"os"

	"code.sphoton.com/sphotonchat/utilities/mmgotool/commands"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
