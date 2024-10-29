package get

import (
	"encoding/json"
	"fmt"
	"os"

	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// passwordCmd represents the password command.
var passwordCmd = &cobra.Command{
	Use:   "password [ID]",
	Short: "Получить полные данные логин-пароля",
	Long:  "Получить полные данные логин-пароля по его ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		data, err := root.Services.GetPassword(id)
		if err != nil {
			printFailed(err)
			os.Exit(1)
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(passwordCmd)
}
