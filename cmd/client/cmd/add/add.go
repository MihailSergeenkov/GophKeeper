package add

import (
	"os"

	root "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

const (
	markFlag        = "mark"
	descriptionFlag = "description"
)

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Загрузить данные",
	Long:  "Загрузить данные на сервер",
}

func init() {
	root.RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringP(markFlag, "m", "", "Пометка для объекта данных")
	addCmd.PersistentFlags().StringP(descriptionFlag, "d", "", "Дополнительное описание для объекта данных")
	if err := addCmd.MarkPersistentFlagRequired(markFlag); err != nil {
		printFailed(root.RootCmd, err)
		os.Exit(1)
	}
}

func printFailed(cmd *cobra.Command, err error) {
	cmd.Printf("Failed: %s", err)
}
