package add

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
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

		if err := services.AddPassword(config.GetConfig(), req); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Add password OK")
	},
}

func init() {
	addCmd.AddCommand(passwordCmd)

	passwordCmd.Flags().StringP("login", "l", "", "Логин для сохранения")
	passwordCmd.Flags().StringP("password", "p", "", "Пароль для сохранения")
	if err := passwordCmd.MarkFlagRequired("login"); err != nil {
		printFailed(err)
		os.Exit(1)
	}
	if err := passwordCmd.MarkFlagRequired("password"); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}
