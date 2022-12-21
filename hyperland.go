package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Hyperland struct {
	sig string
}

func (h Hyperland) listen() error {
	h.sig = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

	ws, err := h.getWorkspaces()
	if err != nil {
		log.Fatal("ws parse", err)
	}

	w := Workspaces{
		Active:     0,
		Workspaces: ws,
	}

	w.toJson()

	//emulating hypr itself which hardcodes the /tmp
	ctl, err := net.Dial("unix", "/tmp/hypr/"+h.sig+"/.socket2.sock")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(ctl)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		parts := strings.Split(msg, ">>")
		switch parts[0] {
		case "workspace":
			ID, _ := strconv.ParseInt(strings.Trim(parts[1], " \n"), 10, 64)
			w.Active = ID
			for i, ws := range w.Workspaces {
				w.Workspaces[i].IsActive = ws.ID == ID
			}
			w.toJson()
		case "destroyworkspace":
			ws, err := h.getWorkspaces()
			if err != nil {
				log.Fatal(err)
			}
			w.Workspaces = ws
			for i, ws := range w.Workspaces {
				w.Workspaces[i].IsActive = ws.ID == w.Active
			}
			w.toJson()
		case "createworkspace":
			ws, err := h.getWorkspaces()
			if err != nil {
				log.Fatal(err)
			}
			w.Workspaces = ws
			for i, ws := range w.Workspaces {
				w.Workspaces[i].IsActive = ws.ID == w.Active
			}
			w.toJson()
		}

	}
}

func (h Hyperland) getWorkspaces() ([]Workspace, error) {
	ws := []Workspace{}

	ctl, err := net.Dial("unix", "/tmp/hypr/"+h.sig+"/.socket.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer ctl.Close()
	// apparently hyperctl internally sends arguments as arg / command
	// https://sourcegraph.com/github.com/hyprwm/Hyprland/-/blob/src/debug/HyprCtl.cpp?L879
	ctl.Write([]byte("j/workspaces"))
	byt, err := io.ReadAll(ctl)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byt, &ws)
	if err != nil {
		return nil, err
	}
	var st SortWorkspaces = ws
	sort.Sort(st)

	return st, nil
}
