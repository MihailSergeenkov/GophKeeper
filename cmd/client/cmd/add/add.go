package add

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
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
	cmd.RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringP(markFlag, "m", "", "Пометка для объекта данных")
	addCmd.PersistentFlags().StringP(descriptionFlag, "d", "", "Дополнительное описание для объекта данных")
	if err := addCmd.MarkPersistentFlagRequired(markFlag); err != nil {
		printFailed(err)
		os.Exit(1)
	}
}

func printFailed(err error) {
	fmt.Printf("Failed: %s", err)
}
