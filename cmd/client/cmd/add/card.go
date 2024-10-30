package add

import (
	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// cardCmd represents the card command.
var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Загрузить данные банковской карты",
	Long:  "Загрузить данные банковской карты на сервер",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		number, _ := cmd.Flags().GetString("number")
		owner, _ := cmd.Flags().GetString("owner")
		expiryDate, _ := cmd.Flags().GetString("expiry-date")
		cvv2, _ := cmd.Flags().GetString("cvv2")
		mark, _ := cmd.Flags().GetString(markFlag)
		description, _ := cmd.Flags().GetString(descriptionFlag)

		req := models.AddCardRequest{
			Number:      number,
			Owner:       owner,
			ExpiryDate:  expiryDate,
			CVV2:        cvv2,
			Mark:        mark,
			Description: description,
		}

		if err := root.Services.AddCard(&req); err != nil {
			printFailed(cmd, err)
			return
		}

		cmd.Println("Add card OK")
	},
}

func init() {
	addCmd.AddCommand(cardCmd)

	cardCmd.Flags().StringP("number", "n", "", "Номер карты для сохранения")
	cardCmd.Flags().StringP("owner", "o", "", "Владелец карты для сохранения")
	cardCmd.Flags().StringP("expiry-date", "e", "", "Дата окончания карты для сохранения")
	cardCmd.Flags().String("cvv2", "", "CVV2 карты для сохранения")
	if err := cardCmd.MarkFlagRequired("number"); err != nil {
		printFailed(addCmd, err)
		return
	}
}
