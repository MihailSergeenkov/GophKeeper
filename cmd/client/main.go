package main //nolint:typecheck // false-positive

import (
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	_ "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/add"
	_ "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/get"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/client/services"
)

func main() {
	cfg := config.GetConfig()
	httpRequests := requests.NewRequests(cfg)
	s := services.Init(cfg, httpRequests)

	cmd.Execute(s)
}
