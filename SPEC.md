# SPEC

## §G GOAL
CLI tool `ptc` converts media (images, GIFs, videos) into Plymouth boot-splash themes. minimal config. zero script writing.

## §C CONSTRAINTS
- Go 1.24+. single static binary.
- Linux-only. Plymouth ! available elsewhere.
- Input: png, jpg, gif, mp4, webm, mov. ? svg.
- Output: standard Plymouth theme dir (`/usr/share/plymouth/themes/<name>/`).
- No runtime deps beyond ffmpeg (video/gif), plymouth/plymouthd (preview).
- Boot-safe assets: ≤ 1920×1080, ≤ 256 colors ? user override.
- Script target: Plymouth Script (JavaScript-like).
- Minimal code comments

## §I INTERFACES
- cmd: `ptc create <name> <media...> [--fps N] [--res WxH] [--loop] [--transition fade|none] [--output-dir path]` → generate theme from media files.
- cmd: `ptc install <dir> [--system-dir path]` → copy to system themes dir. requires root.
- cmd: `ptc preview <dir>` → run `plymouthd --test` with theme. requires root.
- file: `<name>.plymouth` → INI descriptor. keys: `Name`, `Description`, `ModuleName=script`, `ImageDir=assets`, `ScriptFile=name.script`.
- file: `<name>.script` → auto-generated Plymouth Script.
- dir: `assets/` → processed frames/images.

## §V INVARIANTS
V1: ∀ `ptc create` → media files exist & readable before any write.
V2: ∀ `ptc create` → `.plymouth` descriptor parses before any file write.
V3: ∀ `ptc install` → target dir has valid `.plymouth` + `.script` + `assets/` before copy.
V4: ∀ generated `.script` → Plymouth Script API funcs ∈ whitelist. ⊥ undefined func calls.
V5: ∀ asset refs in script → file exists in `assets/` or create fails.
V6: ∀ `ptc preview` → `plymouthd --version` check first. graceful error if missing.
V7: theme name `[a-z0-9_-]+` case-insensitive unique check on install.
V8: ∀ extracted frames → dimensions ≤ config.max_res (default 1920×1080).

## §T TASKS
id|status|task|cites
T1|x|scaffold `cmd/ptc` with cobra, `go.mod`|§C
T2|x|impl media probe (ffprobe wrap): type, dims, duration, frames|§C
T3|x|impl frame extraction: gif→png frames, video→png frames via ffmpeg|V8
T4|x|impl image resize/normalize pipeline (boot-safe)|V8
T5|x|impl `ptc create` with flags → generate .plymouth, .script, assets/|V1,V2,V4,V5
T6|x|impl `.plymouth` parser/validator|V2,V3
T7|x|impl `.script` validator (token-based, API whitelist)|V4
T8|x|impl `ptc install` → validate target + copy to system dir|V3,V6,V7
T9|.|impl `ptc preview` → exec `plymouthd --test`|V6
T10|.|tests for media probe, frame extraction, validators|§V

## §B BUGS
id|date|cause|fix
