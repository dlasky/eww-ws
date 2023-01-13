package main

import (
	"context"
	"log"
	"os"

	"github.com/joshuarubin/go-sway"
)

type Handler struct {
	sway.EventHandler
	s Sway
}

func (h *Handler) WorkspaceEvent(ctx context.Context, evt sway.WindowEvent) {
	var active = 0
	ws, err := h.s.getWorkspaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range ws {
		if w.IsActive {
			active = w.ID
		}
	}

	w := Workspaces{
		Active:     active,
		Workspaces: ws,
	}
	w.toJson()
}

type Sway struct {
	client sway.Client
}

func (s Sway) detect() bool {
	return os.Getenv("SWAYSOCK") != ""
}

func (s Sway) listen() error {
	client, err := sway.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	s.client = client

	var active = 0
	ws, err := s.getWorkspaces()
	if err != nil {
		return err
	}
	for _, w := range ws {
		if w.IsActive {
			active = w.ID
		}
	}

	w := Workspaces{
		Active:     active,
		Workspaces: ws,
	}
	w.toJson()

	h := Handler{
		s: s,
	}

	sway.Subscribe(context.Background(), h, sway.EventTypeWorkspace)
	return nil
}

func (s Sway) getWorkspaces() ([]Workspace, error) {
	ws := []Workspace{}
	swayWs, err := s.client.GetWorkspaces(context.Background())
	for _, w := range swayWs {
		ws = append(ws, Workspace{
			IsActive: w.Focused,
			ID:       int(w.Num),
			Name:     w.Name,
			Monitor:  w.Output,
		})
	}
	return ws, err
}
