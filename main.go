package main

import "log"

type WM interface {
	getWorkspaces() ([]Workspace, error)
	listen() error
}

func main() {
	err := Hyperland{}.listen()
	if err != nil {
		log.Fatal(err)
	}
}
