package cmd

import (
	"gamebooks/pkg/container"
	"github.com/rs/zerolog"
	"os"
)

func Main() {
	lw := zerolog.NewConsoleWriter()
	l := zerolog.New(lw)
	l.Level(zerolog.DebugLevel)

	ctr := container.New()

	if len(os.Args) == 3 {
		if os.Args[1] == "lint" {
			bookName := os.Args[2]
			if err := lintBook(ctr, bookName); err != nil {
				panic(err)
			}
			return
		}

		panic("unknown command: " + os.Args[1])
	}

	http(ctr, l)
}
