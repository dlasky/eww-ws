# EWW-WS

A simple binary that supports updating workspaces as a variable in [EWW](https://github.com/elkowar/eww/).

The goal of this project is to support a consistent JSON format for workspaces in EWW regardless of the window manager. We will also attempt to use IPC interfaces where possible
to avoid comparatively heavy calls to exec like most of the existing 'script' solutions.

Eww-ws relies on the `deflisten` feature of EWW to stream updates as needed. The idea is to avoid polling by using the IPCs of the respective window managers to determine which workspace is active and emit a json object for eww.

Example Eww Configuration:

``` lisp
(deflisten spaces `PATH TO eww-ws`)

(defwidget workspaces []
  (box :class "workspaces"
  (for ws in "${spaces.workspaces}"
    (eventbox
        :width 30
        :class "${ws.is_active ? "active": ''}"
        :onclick "hyprctl dispatch workspace ${ws.id}"
      "${ws.name}"))
  ))

```

Example styling:

```scss
.workspaces {
  color: #928374;
  background-color: #232323;
  font-weight: bold;

  .active {
    color: #89b482;
    border-top: 3px solid #89b482;
    background-color: #141414;
  }
}
```

# Currently Supported:

Single Monitor only -- will be exploring how to support this better in the near future

- [x] Hyprland
- [x] Sway
- [ ] River

# TODO:

- [ ] Multi Monitor
- [ ] WM agnostic 'switch to' workspace command
