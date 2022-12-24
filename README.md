# EWW-WS

A simple binary that supports updating workspaces as a variable in [EWW](https://github.com/elkowar/eww/).

The goal of this project is to support a consistent JSON format for EWW regardless of WM used, to enable simpler overall EWW config. We will also attempt to use IPC interfaces where possible
to avoid comparatively heavy calls to exec.

The app relies on the def listen feature of EWW to stream updates as needed

Example Eww Configuration:

```
(deflisten SPACES :initial "{\"active\":0, \"workspaces\":[]}"
  `~/.config/eww/scripts/eww-ws`)

(defwidget ws []
  (box :class "workspaces"
  (for entry in "${SPACES.workspaces}"
    (button :onclick "hyprctl dispatch workspace ${entry.id}"
      "${entry.name}"))
  ))

```

# Currently Supported:

Single Monitor only -- will be exploring how to support this better in the near future

- [x] Hyprland
- [ ] Sway
- [ ] River