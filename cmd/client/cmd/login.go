package cmd

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Аутентификация пользователя",
	Long:  "Аутентификация пользователя ранее зарегистрированного пользователя, запускается синхронизация",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetString(loginFlag)
		password, _ := cmd.Flags().GetString(passwordFlag)

		req := models.CreateUserTokenRequest{
			Login:    login,
			Password: password,
		}

		if err := services.LoginUser(config.GetConfig(), req); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		if err := services.SyncData(config.GetConfig()); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Login OK")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP(loginFlag, "l", "", "Логин пользователя")
	loginCmd.Flags().StringP(passwordFlag, "p", "", "Пароль пользователя")
	if err := loginCmd.MarkFlagRequired(loginFlag); err != nil {
		printFailed(err)
		os.Exit(1)
	}
	if err := loginCmd.MarkFlagRequired(passwordFlag); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}
