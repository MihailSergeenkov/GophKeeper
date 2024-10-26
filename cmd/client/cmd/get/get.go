package get

import (
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// getCmd represents the add command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить полные данные",
	Long:  "Получить полные данные, конкретного типа",
}

func init() {
	cmd.RootCmd.AddCommand(getCmd)
}
