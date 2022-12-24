package main

import (
	"encoding/json"
	"fmt"
)

type Workspace struct {
	IsActive bool   `json:"is_active"`
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Monitor  string `json:"monitor,omitempty"`
}

type Workspaces struct {
	Active     int         `json:"active,omitempty"`
	Workspaces []Workspace `json:"workspaces,omitempty"`
}

type SortWorkspaces []Workspace

func (s SortWorkspaces) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortWorkspaces) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
func (s SortWorkspaces) Len() int {
	return len(s)
}

func (ws Workspaces) toJson() {
	byt, _ := json.Marshal(ws)
	fmt.Println(string(byt))
}
