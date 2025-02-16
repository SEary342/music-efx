<img src="logoBig.png" alt="MusicEFX Logo" width="200" height="200">

# MusicEFX

A terminal-based MP3 player built with Python. MusicEFX allows users to navigate directories, select MP3 files, and play them directly from the terminal. With a simple TUI interface, the app offers playlist management, folder navigation, search functionality, and playback controls.

## Features

- **Main Menu Navigation**: Navigate through different sections such as Auto-Playlist, Playlists, Folder Navigation, and Search.
- **Folder Navigation**: Browse through directories to find MP3 files.
- **Playlist Management**: Access existing playlists and play MP3 files from them.
- **Playback Control**: Select an MP3 file to play, stop playback, and navigate through options.
- **Folder/Directory** Navigation: Navigate nested directories to select and play MP3 files.
  
## Prerequisites

Before running the app, you will need:
- A terminal environment for interacting with the TUI (Text User Interface).

## Usage

### Running the Application

Download the release from GitHub and run the app by specifying a directory path containing MP3 files:

**Linux:**
```bash
./music-efx
```

#### Stopping Playback

To stop the playback of a song, press `s`.

## Source Installation

### Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/seary342/music-efx.git
cd music-efx
```

### Install Dependencies

Make sure you have [uv](https://docs.astral.sh/uv/) available and run the following command

```bash
uv run main.py
```

Alternatively, you can pip install and run the main.py:

```bash
pip install .
python main.py
```

### Build the Application

To build the application, use the following commands:

```bash
uv build

# Set the applicable environment variables: https://ofek.dev/pyapp/latest/examples/#custom-embedded-local-distribution

# Follow the pyapp instructions: https://ofek.dev/pyapp/latest/how-to/
```

## License

This project is open source and available under the [MIT License](LICENSE).

