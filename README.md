# EWW-WS

A simple binary that supports updating workspaces as a variable in [EWW](https://github.com/elkowar/eww/).

The goal of this project is to support a consistent JSON format for EWW regardless of WM used, to enable simpler overall EWW config. We will also attempt to use IPC interfaces where possible
to avoid comparatively heavy calls to exec like most of the existing 'script' solutions.

The tool relies on the `def listen` feature of EWW to stream updates as needed. The idea being that we avoid polling by utilizing the IPCs of the respective window managers to determine which workspace is active, then we emit an update json object which eww picks up.

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
- [x] Sway
- [ ] River

# TODO:

- [ ] Multi Monitor
- [ ] WM agnostic 'switch to' workspace command