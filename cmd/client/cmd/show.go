package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// showCmd represents the show command.
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Показать сохранненную информацию",
	Long:  "Показать сохранненную информацию для дальнейшей загрузки",
	Run: func(cmd *cobra.Command, args []string) {
		sync, _ := cmd.Flags().GetBool("sync")

		if sync {
			if err := services.SyncData(config.GetConfig()); err != nil {
				printFailed(err)
				os.Exit(1)
			}
		}

		userData := services.GetData(config.GetConfig())

		b, err := json.MarshalIndent(userData, "", "  ")
		if err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println(string(b))
	},
}

func init() {
	RootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolP("sync", "s", false, "Обновить синхронизацию базовой информацииы")
}
