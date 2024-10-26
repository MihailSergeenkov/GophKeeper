package add

import (
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Загрузить данные",
	Long:  "Загрузить данные на сервер",
}

func init() {
	cmd.RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringP("mark", "m", "", "Пометка для объекта данных")
	addCmd.MarkPersistentFlagRequired("mark")
	addCmd.PersistentFlags().StringP("description", "d", "", "Дополнительное описание для объекта данных")
}
