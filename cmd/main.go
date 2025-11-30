package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{}
	rootCommand.AddCommand(CreateBoilerplate())
	rootCommand.Execute()
}
