package main

import (
	"fmt"
	"os"

	"boxed/cmd"
	"boxed/internal/render"
)

func main() {
	renderer := render.NewLipGlossRenderer()
	executor := cmd.NewExecutor(renderer, os.Stdout)
	rootCmd := cmd.NewRootCmd(executor)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
