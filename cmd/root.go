package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/damonto/ms365/internal/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "ms365",
	Short: "Microsoft 365 RESTful API",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// configuration file path
var cfgPath string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgPath, "conf", "configs/config.toml", "The application configuration file path")
}

func initConfig() {
	config.ReadConfig(cfgPath)
}
