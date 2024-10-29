package cmd

import (
	"fmt"
	"os"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

const (
	loginFlag    = "login"
	passwordFlag = "password"
)

var version = "0.0.1"
var cfgFile string

// Servicer интерфейс для сервисов клиента.
type Servicer interface {
	RegisterUser(req models.RegisterUserRequest) error
	LoginUser(req models.CreateUserTokenRequest) error
	SyncData() error
	LogoutUser() error
	GetData() []models.UserData
	AddPassword(req models.AddPasswordRequest) error
	GetPassword(id string) (models.Password, error)
	AddCard(req *models.AddCardRequest) error
	GetCard(id string) (models.Card, error)
	AddText(req models.AddTextRequest) error
	GetText(id string) (models.Text, error)
	AddFile(filePath, mark, description string) error
	GetFile(id, dir string) error
}

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "client",
	Short: "Клиент для GophKeeper",
	Long: `Клиент для GophKeeper позволяющий добавлять и просматривать:
	- логин-пароли;
	- банковские карты;
	- тексты;
	- файлы`,
	Version: version,
}

var Services Servicer

func Execute(s Servicer) {
	Services = s
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Initializer(&cfgFile))
	RootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "config file (default is $HOME/.goph-keeeper.yaml)")
}

func printFailed(err error) {
	fmt.Printf("Failed: %s", err)
}
