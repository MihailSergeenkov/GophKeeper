package add

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

// cardCmd represents the card command
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
		mark, _ := cmd.Flags().GetString("mark")
		description, _ := cmd.Flags().GetString("description")

		req := models.AddCardRequest{
			Number:      number,
			Owner:       owner,
			ExpiryDate:  expiryDate,
			CVV2:        cvv2,
			Mark:        mark,
			Description: description,
		}

		if err := services.AddCard(config.GetConfig(), req); err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		fmt.Println("Add text OK")
	},
}

func init() {
	addCmd.AddCommand(cardCmd)

	cardCmd.Flags().StringP("number", "n", "", "Номер карты для сохранения")
	cardCmd.MarkFlagRequired("number")
	cardCmd.Flags().StringP("owner", "o", "", "Владелец карты для сохранения")
	cardCmd.Flags().StringP("expiry-date", "e", "", "Дата окончания карты для сохранения")
	cardCmd.Flags().String("cvv2", "", "CVV2 карты для сохранения")
}
