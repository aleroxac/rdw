# RDW (Remind Drink Water)

A lightweight desktop water tracking application for Linux, built with Go.

## Features
- **System Tray Integration**: Stealthy operation in the background.
- **Periodic Notifications**: Customizable alerts to keep you hydrated.
- **Daily Tracking**: Log your intake and monitor progress against goals.
- **Desktop Native**: Built specifically for Linux environments (optimized for i3wm/Sway).

## Prerequisites
- Go 1.21+
- `libayatana-appindicator-glib`
- `gtk3`

## Installation
```bash
# ArchLinux
yay -S --needed libayatana-appindicator-glib gtk3
go mod download
go build -o rdw .
sudo mv rdw /usr/local/bin
```

## Usage
``` bash
cat << EOF > ${HOME}/.config/systemd/user/rdw.service
[Unit]
Description=RDW - Water Tracker Tray App
After=graphical-session.target

[Service]
ExecStart=/usr/local/bin/rdw
Restart=on-failure
# Garante que ele consiga falar com o X11/Wayland
Environment=DISPLAY=:0
Environment=XAUTHORITY=%h/.Xauthority
Environment=XDG_CURRENT_DESKTOP=Unity

[Install]
WantedBy=graphical-session.target
EOF

systemctl --user daemon-reload
systemctl --user enable rdw
systemctl --user start rdw
systemctl --user status rdw
```
