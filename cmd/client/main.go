package main

import (
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	_ "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/add"
	_ "github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/get"
)

func main() {
	cmd.Execute()
}
