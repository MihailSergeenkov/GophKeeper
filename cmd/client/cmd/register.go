package cmd

import (
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
			printFailed(cmd, err)
			return
		}

		cmd.Println("Register OK")
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP(loginFlag, "l", "", "Логин пользователя")
	registerCmd.Flags().StringP(passwordFlag, "p", "", "Пароль пользователя")
	_ = registerCmd.MarkFlagRequired(loginFlag)
	_ = registerCmd.MarkFlagRequired(passwordFlag)
}
