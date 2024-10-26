package cmd

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрация пользователя",
	Long:  "Регистрация нового пользователя сервиса, после успешной регистрации необходимо выполнить login",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetString("login")
		password, _ := cmd.Flags().GetString("password")

		req := models.RegisterUserRequest{
			Login:    login,
			Password: password,
		}

		err := services.RegisterUser(config.GetConfig(), req)
		if err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		fmt.Println("Register OK")
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("login", "l", "", "Логин пользователя")
	registerCmd.MarkFlagRequired("login")
	registerCmd.Flags().StringP("password", "p", "", "Пароль пользователя")
	registerCmd.MarkFlagRequired("password")
}
