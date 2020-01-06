package main

import (
	controller "github.com/micro-community/x-micro-edge/end"
	_ "github.com/micro/go-micro"
)

func start() {
	go controller.RunProc()
}
