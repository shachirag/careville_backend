package main

import "careville_backend/app"

// @title Care Ville
// @version 0.0.1
// @description Care Ville Backend (in GoLang)
// @contact.name Chirag Sharma
// @license.name MIT
// @host zn2j5663-5065.inc1.devtunnels.ms
// @BasePath /
func main() {
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}