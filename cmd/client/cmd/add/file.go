package add

import (
	"fmt"
	"os"

	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command.
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Загрузить файл",
	Long:  "агрузить файл на сервер",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		mark, _ := cmd.Flags().GetString(markFlag)
		description, _ := cmd.Flags().GetString(descriptionFlag)

		if err := root.Services.AddFile(file, mark, description); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Add file OK")
	},
}

func init() {
	addCmd.AddCommand(fileCmd)

	fileCmd.Flags().StringP("file", "f", "", "Файл для сохранения")
	if err := fileCmd.MarkFlagRequired("file"); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}
