package cmd

import (
	"fmt"

	"github.com/aguidirh/go-rag-chatbot/internal/cmd"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "go-rag-chatbot",
		Short: "Provides RAG capabilities",
	}
)

func init() {
	rootCmd.PersistentFlags().String("config-path", "~/.go-rag-chatbot", "Path to the configuration file")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level")
	rootCmd.AddCommand(cmd.ServeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("error: %v", err))
	}
}
