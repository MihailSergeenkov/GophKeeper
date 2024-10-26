package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
	"github.com/spf13/cobra"
)

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text [ID]",
	Short: "Получить текст",
	Long:  "Получить текст по его ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		data, err := services.GetText(config.GetConfig(), id)
		if err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Failed: %s", err)
			os.Exit(1)
		}

		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(textCmd)
}
