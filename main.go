package main

import "github.com/8micro/gounits/system"

import (
	"log"
	"os"
)

func main() {

	server, err := NewWechatServer()
	if err != nil {
		log.Printf("server error, %s\n", err.Error())
		os.Exit(system.ErrorExitCode(err))
	}

	defer func() {
		exitCode := 0
		if err := server.Stop(); err != nil {
			log.Printf("server stop error, %s\n", err.Error())
			exitCode = system.ErrorExitCode(err)
		}
		os.Exit(exitCode)
	}()

	if err = server.Startup(); err != nil {
		log.Printf("server startup error, %s\n", err.Error())
		os.Exit(system.ErrorExitCode(err))
	}
	system.InitSignal(nil)
}
