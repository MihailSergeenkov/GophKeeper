package cmd

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация пользователя",
	Long:  "Регистрация нового пользователя сервиса, после успешной регистрации необходимо выполнить login",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetString(loginFlag)
		password, _ := cmd.Flags().GetString(passwordFlag)

		req := models.RegisterUserRequest{
			Login:    login,
			Password: password,
		}

		if err := Services.RegisterUser(req); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Register OK")
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP(loginFlag, "l", "", "Логин пользователя")
	registerCmd.Flags().StringP(passwordFlag, "p", "", "Пароль пользователя")
	if err := registerCmd.MarkFlagRequired(loginFlag); err != nil {
		printFailed(err)
		os.Exit(1)
	}
	if err := registerCmd.MarkFlagRequired(passwordFlag); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}
