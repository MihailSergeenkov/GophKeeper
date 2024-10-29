package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command.
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Разлогин пользователя",
	Long:  "Разлогин пользователя, удаление синхронизированной базовой информации",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Services.LogoutUser(); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Logout OK")
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
