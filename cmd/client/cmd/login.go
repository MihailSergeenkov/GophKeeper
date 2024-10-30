package cmd

import (
	"os"

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

		if err := Services.LoginUser(req); err != nil {
			printFailed(cmd, err)
			return
		}

		if err := Services.SyncData(); err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Println("Login OK")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP(loginFlag, "l", "", "Логин пользователя")
	loginCmd.Flags().StringP(passwordFlag, "p", "", "Пароль пользователя")
	if err := loginCmd.MarkFlagRequired(loginFlag); err != nil {
		printFailed(RootCmd, err)
		os.Exit(1)
	}
	if err := loginCmd.MarkFlagRequired(passwordFlag); err != nil {
		printFailed(RootCmd, err)
		os.Exit(1)
	}
}
