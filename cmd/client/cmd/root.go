package cmd

import (
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/spf13/cobra"
)

var version = "0.0.1"
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "client",
	Short: "Клиент для GophKeeper",
	Long: `Клиент для GophKeeper позволяющий добавлять и просматривать:
	- логин-пароли;
	- банковские карты;
	- тексты;
	- файлы`,
	Version: version,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Initializer(&cfgFile))
	RootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "config file (default is $HOME/.goph-keeeper.yaml)")
}
