# MusicEFX

A terminal-based MP3 player built with Go. This app lets users navigate through a list of MP3 files, select a file, and play it directly from the terminal interface. It's a simple, text-based MP3 player, featuring basic playback controls such as stop and navigation through files.

**Work in Progress**: This project is under active development. Some features may not be fully implemented, and the user experience may change as the app evolves.

## Features

- **File Navigation**: Navigate through MP3 files using the up and down arrow keys.
- **Playback Control**: Select a file to play with the `Enter` key.
- **Stop Playback**: Stop the current playback using the `Esc` key or `Q`.
- **Pagination**: MP3 files are displayed in pages, showing up to 10 files per page.
  
## Prerequisites

Before running the app, you will need:

- Go 1.18 or higher.
- A terminal environment for interacting with the TUI (Text User Interface).

## Installation

### Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/seary342/music-efx.git
cd music-efx
```

### Install Dependencies

Run the following command to install necessary Go dependencies:

```bash
go mod tidy
```

### Build the Application

To build the application, use the following command:

```bash
go build -o MusicEFX .
```

This will create the executable `MusicEFX`.

## Usage

### Running the Application

Run the app by specifying a directory path containing MP3 files:

```bash
./MusicEFX /path/to/mp3/files
```

If no directory is specified, the app will prompt you to input a directory path.

#### Controls

- **Up/Down Arrow Keys:** Navigate through the list of MP3 files.
- **Enter:** Play the selected MP3 file.
- **Esc:** Stop the current playback (if no song is playing, this exits the app).
- **Q:** Stop the current playback and quit the application.

#### Stopping Playback

To stop the playback of a song, press `Esc` or `Q`.

## Development

This project is still in development, and several features and enhancements are planned for future releases, including:

    Enhanced error handling and feedback.
    User interface improvements.
    Playlist management features.
    Auto-Playlist support with timing systems

Feel free to contribute by submitting issues or pull requests!

## License

This project is open source and available under the MIT License.

