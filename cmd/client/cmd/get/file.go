package get

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file [ID]",
	Short: "Получить файл",
	Long:  "Получить файл по его ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		dir, _ := cmd.Flags().GetString("upload-dir")
		if err := services.GetFile(config.GetConfig(), id, dir); err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		fmt.Printf("File load in %s", dir)
	},
}

func init() {
	getCmd.AddCommand(fileCmd)

	fileCmd.Flags().StringP("upload-dir", "u", ".", "upload file directory (default is ${pwd})")
}
