package get

import (
	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command.
var fileCmd = &cobra.Command{
	Use:   "file [ID]",
	Short: "Получить файл",
	Long:  "Получить файл по его ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		dir, _ := cmd.Flags().GetString("upload-dir")
		if err := root.Services.GetFile(id, dir); err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Printf("File load in %s", dir)
	},
}

func init() {
	getCmd.AddCommand(fileCmd)

	fileCmd.Flags().StringP("upload-dir", "u", ".", "upload file directory (default is ${pwd})")
}
