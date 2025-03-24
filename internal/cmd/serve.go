package cmd

import (
	"fmt"
	"strings"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/httpserver"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves an API which enables the configuration and querying of RAG sources",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.InheritedFlags().GetString("config-path")
		if err != nil {
			panic(fmt.Errorf("unable to parse config-path argument. %v", err))
		}
		logLevel, err := cmd.InheritedFlags().GetString("log-level")
		if err != nil {
			panic(fmt.Errorf("unable to parse log-level argument. %v", err))
		}

		logger := logrus.New()

		parsedLevel, err := logrus.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			panic(fmt.Errorf("unable to parse log level. %v", err))
		}
		logger.SetLevel(parsedLevel)

		server := httpserver.HttpServer{
			ConfigPath: configPath,
			Log:        logger,
		}
		server.Run()
	},
}
