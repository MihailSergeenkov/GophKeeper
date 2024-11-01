package add

import (
	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// passwordCmd represents the password command.
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Загрузить данные логин-пароля",
	Long:  "Загрузить данные логин-пароля на сервер",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetString("login")
		password, _ := cmd.Flags().GetString("password")
		mark, _ := cmd.Flags().GetString(markFlag)
		description, _ := cmd.Flags().GetString(descriptionFlag)

		req := models.AddPasswordRequest{
			Login:       login,
			Password:    password,
			Mark:        mark,
			Description: description,
		}

		if err := root.Services.AddPassword(req); err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Println("Add password OK")
	},
}

func init() {
	addCmd.AddCommand(passwordCmd)

	passwordCmd.Flags().StringP("login", "l", "", "Логин для сохранения")
	passwordCmd.Flags().StringP("password", "p", "", "Пароль для сохранения")
	_ = passwordCmd.MarkFlagRequired("login")
	_ = passwordCmd.MarkFlagRequired("password")
}
