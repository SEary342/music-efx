# MusicEFX

MusicEFX is a terminal-based MP3 player built with Python and Textual. It lets you browse music folders, build a playlist from a selected directory, and control playback from a keyboard-driven TUI.

## Features

- Browse local folders and pick audio files
- Build a playlist from a selected folder
- Play, pause, stop, previous/next, random, and reroll controls
- Keyboard shortcuts for quick control: Space (play/pause), s (stop), p (previous), n (next), and r (toggle random)

## Requirements

- Python 3.14+
- [uv](https://docs.astral.sh/uv/)

## Development setup

```bash
git clone https://github.com/seary342/music-efx.git
cd music-efx
uv sync --group dev
```

## Run the app

From the repository root, start the TUI with:

```bash
uv run src/music_efx
```

## Development commands

The project uses Taskipy for common workflows:

```bash
uv run task test
uv run task lint
uv run task format
uv run task build
```

Useful maintenance commands:

```bash
uv run task check
uv run task clean
```

## Build

To build a standalone application with PyInstaller, run:

```bash
uv run task build
```

The packaged binary is written to the dist directory.

## License

This project is open source and available under the [MIT License](LICENSE).

