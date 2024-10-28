package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// cardCmd represents the card command.
var cardCmd = &cobra.Command{
	Use:   "card [ID]",
	Short: "Получить полные данные банковской карты",
	Long:  "Получить полные данные банковской карты по его ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		data, err := services.GetCard(config.GetConfig(), id)
		if err != nil {
			printFailed(err)
			os.Exit(1)
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			printFailed(err)
			os.Exit(1)
		}

		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(cardCmd)
}
