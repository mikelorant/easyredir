package main

import (
	"fmt"
	"os"

	"github.com/mikelorant/easyredir-cli/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "error: %v\n", err)
		os.Exit(1)
	}
}
