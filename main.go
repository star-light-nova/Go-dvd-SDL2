package main

import (
	"dvd/app"
	"fmt"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}
