package cmd

import (
	"fmt"
	"os"

	_ "github.com/mattes/migrate/database/postgres" // driver
	_ "github.com/mattes/migrate/source/file"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/caquillo07/graphql-server-demo/conf"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gql",
	Short: "Test GraphQL Server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() { conf.InitViper(cfgFile) })
	cobra.OnInitialize(initLogging)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().Bool("dev-log", false, "Development logging")
}

// initConfig reads in config file and ENV variables if set.
func initLogging() {
	var logger *zap.Logger
	if val, _ := rootCmd.PersistentFlags().GetBool("dev-log"); val {
		logger, _ = zap.NewDevelopment()
		logger.Info("Development logging enabled")
	} else {
		logger, _ = zap.NewProduction()
	}

	logger.Info("GraphQL Server Started")
	zap.ReplaceGlobals(logger)
}
