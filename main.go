package main

import "log"

type WM interface {
	detect() bool
	getWorkspaces() ([]Workspace, error)
	listen() error
}

var managers = []WM{Hyperland{}, Sway{}}

func main() {

	for _, manager := range managers {
		if manager.detect() {
			err := manager.listen()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
