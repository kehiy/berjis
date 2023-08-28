package main

import (
	"fmt"

	"github.com/kehiy/berjis"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "berjis",
		Short: "Berjis daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("use --help for more info")
		},
		Version: berjis.Agent(),
	}

	buildRunCommand(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
