# ptc — Plymouth Theme Creator

CLI tool converting images, GIFs, and videos into Plymouth boot-splash themes. Single static binary. Zero script writing.

## Requirements

- Go 1.24+ (build from source)
- Linux (Plymouth only runs on Linux)
- `ffmpeg` (video/GIF processing)
- `plymouthd` (preview only)

## Install

```bash
go install github.com/dnevb/go-ptc/cmd/ptc@latest
```

Or build from source:

```bash
go build -o ~/.local/bin/ptc ./cmd/ptc
```

## Usage

### Create theme

```bash
ptc create mytheme video.mp4 --fps 30 --res 1920x1080 --loop
ptc create mytheme image.png --res 800x600
```

Generates `mytheme/` with `.plymouth`, `.script`, and `assets/`.

### Install theme

```bash
sudo ptc install mytheme
```

Copies to `/usr/share/plymouth/themes/`. Requires root.

### Preview theme

```bash
sudo ptc preview mytheme
```

Runs `plymouthd` + `plymouth --show-splash` for 5 seconds.
