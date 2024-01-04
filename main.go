package main

import "careville_backend/app"

// @title Care Ville
// @version 0.0.1
// @description Care Ville Backend (in GoLang)
// @contact.name Chirag Sharma
// @license.name MIT
// @host mmpn2atpcm.eu-west-1.awsapprunner.com
// @BasePath /
func main() {
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}

// mmpn2atpcm.eu-west-1.awsapprunner.com