package cmd

import (
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Разлогин пользователя",
	Long:  "Разлогин пользователя, удаление синхронизированной базовой информации",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		services.LogoutUser(config.GetConfig())
		fmt.Println("Logout OK")
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
