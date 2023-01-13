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

func (h Hyperland) detect() bool {
	return os.Getenv("HYPERLAND_INSTANCE_SIGNATURE") != ""
}

func (h Hyperland) listen() error {
	h.sig = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

	ws, err := h.getWorkspaces()
	if err != nil {
		log.Fatal("ws parse", err)
	}

	active, err := h.getActiveWorkspace()
	if err != nil {
		return err
	}

	w := Workspaces{
		Active:     active,
		Workspaces: ws,
	}
	for i, ws := range w.Workspaces {
		w.Workspaces[i].IsActive = ws.ID == active
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
			w.Active = int(ID)
			for i, ws := range w.Workspaces {
				w.Workspaces[i].IsActive = ws.ID == int(ID)
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

type HyperlandMonitors struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	RefreshRate     float64 `json:"refreshRate"`
	X               int     `json:"x"`
	Y               int     `json:"y"`
	ActiveWorkspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"activeWorkspace"`
	Reserved   []int   `json:"reserved"`
	Scale      float64 `json:"scale"`
	Transform  int     `json:"transform"`
	Focused    bool    `json:"focused"`
	DpmsStatus bool    `json:"dpmsStatus"`
}

func (h Hyperland) getActiveWorkspace() (int, error) {
	ctl, err := net.Dial("unix", "/tmp/hypr/"+h.sig+"/.socket.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer ctl.Close()
	// apparently hyperctl internally sends arguments as arg / command
	// https://sourcegraph.com/github.com/hyprwm/Hyprland/-/blob/src/debug/HyprCtl.cpp?L879
	ctl.Write([]byte("j/monitors"))
	byt, err := io.ReadAll(ctl)
	if err != nil {
		return 0, err
	}
	monitors := []HyperlandMonitors{}
	err = json.Unmarshal(byt, &monitors)
	if err != nil {
		return 0, err
	}
	return monitors[0].ActiveWorkspace.ID, nil
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
