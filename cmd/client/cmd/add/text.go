package add

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
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

		if err := services.AddText(config.GetConfig(), req); err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println("Add text OK")
	},
}

func init() {
	addCmd.AddCommand(textCmd)

	textCmd.Flags().StringP("text", "t", "", "Текст для сохранения")
	if err := textCmd.MarkFlagRequired("text"); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}
