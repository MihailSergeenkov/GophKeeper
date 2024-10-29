package get

import (
	"encoding/json"
	"fmt"
	"os"

	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
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
		data, err := root.Services.GetCard(id)
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
