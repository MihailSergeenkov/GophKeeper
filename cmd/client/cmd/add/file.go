package add

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Загрузить файл",
	Long:  "агрузить файл на сервер",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		mark, _ := cmd.Flags().GetString("mark")
		description, _ := cmd.Flags().GetString("description")

		if err := services.AddFile(config.GetConfig(), file, mark, description); err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		fmt.Println("Add file OK")
	},
}

func init() {
	addCmd.AddCommand(fileCmd)

	fileCmd.Flags().StringP("file", "f", "", "Файл для сохранения")
	fileCmd.MarkFlagRequired("file")
}
