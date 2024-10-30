package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command.
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Показать сохранненную информацию",
	Long:  "Показать сохранненную информацию для дальнейшей загрузки",
	Run: func(cmd *cobra.Command, args []string) {
		sync, _ := cmd.Flags().GetBool("sync")

		fmt.Print("wrsdfsdfsdf", sync)
		if sync {
			if err := Services.SyncData(); err != nil {
				printFailed(cmd, err)
				return
			}
		}

		userData := Services.GetData()
		b, err := json.MarshalIndent(userData, "", "  ")
		if err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Println(string(b))
	},
}

func init() {
	RootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolP("sync", "s", false, "Обновить синхронизацию базовой информацииы")
}
