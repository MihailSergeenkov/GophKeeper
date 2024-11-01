package get

import (
	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// getCmd represents the add command.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить полные данные",
	Long:  "Получить полные данные, конкретного типа",
}

func init() {
	root.RootCmd.AddCommand(getCmd)
}

func printFailed(cmd *cobra.Command, err error) {
	cmd.Printf("Failed: %s", err)
}
