package add

import (
	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// textCmd represents the text command.
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Загрузить текстовые данные",
	Long:  "Загрузить текстовые данные на сервер",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		text, _ := cmd.Flags().GetString("text")
		mark, _ := cmd.Flags().GetString(markFlag)
		description, _ := cmd.Flags().GetString(descriptionFlag)

		req := models.AddTextRequest{
			Data:        text,
			Mark:        mark,
			Description: description,
		}

		if err := root.Services.AddText(req); err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Println("Add text OK")
	},
}

func init() {
	addCmd.AddCommand(textCmd)

	textCmd.Flags().StringP("text", "t", "", "Текст для сохранения")
	if err := textCmd.MarkFlagRequired("text"); err != nil {
		printFailed(addCmd, err)
		return
	}
}
