package cmd

import (
	"os"

	"github.com/damonto/msonline/internal/pkg/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msonline",
	Short: "Microsoft Online RESTful API",
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
