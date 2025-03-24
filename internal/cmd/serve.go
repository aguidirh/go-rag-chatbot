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
		var configPath, logLevel, vectordbHost, listenHost string
		var vectordbPort, listenPort int
		var err error

		if cmd.Flags().Changed("config-path") {
			configPath, err = cmd.Flags().GetString("config-path")
			if err != nil {
				panic(fmt.Errorf("unable to parse config-path argument. %v", err))
			}
		}

		if cmd.Flags().Changed("log-level") {
			logLevel, err = cmd.Flags().GetString("log-level")
			if err != nil {
				panic(fmt.Errorf("unable to parse log-level argument. %v", err))
			}
		}

		if cmd.Flags().Changed("vectordb-host") {
			vectordbHost, err = cmd.Flags().GetString("vectordb-host")
			if err != nil {
				panic(fmt.Errorf("unable to parse vectordb-host argument. %v", err))
			}
		}

		if cmd.Flags().Changed("vectordb-port") {
			vectordbPort, err = cmd.Flags().GetInt("vectordb-port")
			if err != nil {
				panic(fmt.Errorf("unable to parse vectordb-port argument. %v", err))
			}
		}

		if cmd.Flags().Changed("listen-address") {
			listenHost, err = cmd.Flags().GetString("listen-address")
			if err != nil {
				panic(fmt.Errorf("unable to parse listen-address argument. %v", err))
			}
		}

		if cmd.Flags().Changed("listen-port") {
			listenPort, err = cmd.Flags().GetInt("listen-port")
			if err != nil {
				panic(fmt.Errorf("unable to parse listen-port argument. %v", err))
			}
		}

		logger := logrus.New()
		if len(logLevel) == 0 {
			logLevel = "info"
		}
		parsedLevel, err := logrus.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			panic(fmt.Errorf("unable to parse log level. %v", err))
		}
		logger.SetLevel(parsedLevel)

		server := httpserver.HttpServer{
			ConfigPath:   configPath,
			Log:          logger,
			BindAddress:  listenHost,
			BindPort:     listenPort,
			VectorDBHost: vectordbHost,
			VectorDBPort: vectordbPort,
		}
		server.Run()
	},
}

func init() {
	ServeCmd.Flags().String("vectordb-host", "127.0.0.1", "address of vector DB service")
	ServeCmd.Flags().Int("vectordb-port", 6333, "port of vector DB service")
	ServeCmd.Flags().Int("listen-port", 8080, "port on which to listen")
	ServeCmd.Flags().String("listen-address", "127.0.0.1", "address on which to listen")
}
