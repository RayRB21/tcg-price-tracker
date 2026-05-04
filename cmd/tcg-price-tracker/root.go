package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tcg",
	Short: "Riftbound TCG price tracker",
	Long:  "Tracks Riftbound card prices over time and detects trends and spikes.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
