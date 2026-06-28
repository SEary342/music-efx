import importlib.metadata
from pathlib import Path

from textual.app import App, ComposeResult
from textual.containers import Horizontal, Vertical
from textual.widgets import Button, Footer, Header, ProgressBar, Static, Tree

from src.data import MusicEFXData
from src.player import Playlist

from .widgets import FileSystemTree


class MusicEFX(App):
    """TUI Music Player Application."""

    CSS = """
        ProgressBar {
            margin-right: 2;
        }

        #folder_tree {
            height: 70%;
        }

        #controls {
            align-vertical: bottom;
        }

        Button {
            margin-right: 1;
        }
        """

    def __init__(self, root_folder: Path, **kwargs):
        super().__init__(**kwargs)
        self.title = f"MusicEFX {importlib.metadata.version('music-efx')}"
        self.root_folder = root_folder
        self.playlist = Playlist()
        self._poll_timer = None
        self._manually_stopped = False

    def compose(self) -> ComposeResult:
        """Compose the layout of the application."""
        data_provider = MusicEFXData(self.root_folder)
        yield Header()
        yield FileSystemTree(
            label="Music Folders", data_provider=data_provider, id="folder_tree"
        )
        with Vertical(id="info"):
            yield Static("\nNo song playing", id="song_status")
            with Horizontal():
                yield ProgressBar(total=100, id="progress_bar", show_eta=False)
                yield Static("", id="eta_display")
        with Horizontal(id="controls"):
            yield Button(
                "Play/Pause", id="play_pause_button", disabled=True, classes="btn"
            )
            yield Button("Stop", id="stop_button", disabled=True, classes="btn")
            yield Button("Previous", id="prev_button", disabled=True, classes="btn")
            yield Button("Next", id="next_button", disabled=True, classes="btn")
            yield Button("Random", id="random_button", disabled=True, classes="btn")
            yield Button("Re-roll", id="reroll_button", disabled=True, classes="btn")
        yield Footer()

    def on_mount(self) -> None:
        # Start a periodic timer to update the UI (especially button states)
        self._poll_timer = self.set_interval(0.5, self._update_ui)

    def on_tree_node_selected(self, event: Tree.NodeSelected) -> None:
        """Handle folder or file selection in the tree."""
        if not event.node.data:
            return

        selected_path = Path(event.node.data)

        if selected_path.is_dir():
            tracks = [
                f
                for f in sorted(selected_path.glob("*"))
                if f.is_file() and f.suffix.lower() in [".mp3", ".wav", ".ogg"]
            ]
            if tracks:
                self.playlist.set_playlist(tracks)
                self.playlist.player.stop()
                self._manually_stopped = True
        elif selected_path.suffix.lower() in [".mp3", ".wav", ".ogg"]:
            tracks = [
                f
                for f in sorted(selected_path.parent.glob("*"))
                if f.is_file() and f.suffix.lower() in [".mp3", ".wav", ".ogg"]
            ]
            self.playlist.set_playlist(tracks)
            try:
                idx = tracks.index(selected_path)
            except ValueError:
                idx = 0
            self.playlist.play_playlist(idx)
            self._manually_stopped = False

        self._update_ui()

    def _update_ui(self) -> None:
        """Update buttons and labels based on player status."""
        song_status = self.query_one("#song_status", Static)
        progress_bar = self.query_one("#progress_bar", ProgressBar)
        eta_display = self.query_one("#eta_display", Static)

        status_text = "\nNo song playing"
        if self.playlist.player.is_playing:
            name = (
                MusicEFXData.format_song_name(self.playlist.player.current_song)
                if self.playlist.player.current_song
                else "Unknown"
            )
            status_text = f"\nPlaying: {name}"
        elif self.playlist.player.is_paused:
            name = (
                MusicEFXData.format_song_name(self.playlist.player.current_song)
                if self.playlist.player.current_song
                else "Unknown"
            )
            status_text = f"\nPaused: {name}"
        else:
            # If a song finished playing naturally, auto-advance
            if (
                self.playlist.player.current_song
                and not self.playlist.player.active
                and not self.playlist.player.is_paused
                and not self._manually_stopped
            ):
                self.playlist.next_song()
                self._manually_stopped = False

        # Append next song display
        next_song = self.playlist.get_next_song()
        if next_song:
            next_name = MusicEFXData.format_song_name(next_song)
            status_text += f" => Next: {next_name}"

        song_status.update(status_text)

        # Update Progress Bar and ETA
        if self.playlist.player.active or self.playlist.player.is_paused:
            duration = self.playlist.player.get_duration()
            progress = self.playlist.player.get_progress()
            if duration > 0:
                progress_bar.progress = (progress / duration) * 100
                remaining = max(0, duration - progress)
                minutes, seconds = divmod(int(remaining), 60)
                eta_display.update(f"{minutes:02}:{seconds:02}")
            else:
                progress_bar.progress = 0
                eta_display.update("")
        else:
            progress_bar.progress = 0
            eta_display.update("")

        # Update button disabled states
        has_playlist = len(self.playlist.playlist) > 0
        for btn in self.query(".btn").results(Button):
            if btn.id == "prev_button":
                btn.disabled = not has_playlist or len(self.playlist.played_items) <= 1
            elif btn.id == "reroll_button":
                btn.disabled = not has_playlist
                btn.display = self.playlist.random_mode
            else:
                btn.disabled = not has_playlist

    def on_button_pressed(self, event: Button.Pressed) -> None:
        """Handle button presses."""
        if event.button.id == "play_pause_button":
            if self.playlist.player.is_playing:
                self.playlist.player.toggle_pause()
            elif self.playlist.player.is_paused:
                self.playlist.player.toggle_pause()
                self._manually_stopped = False
            else:
                self.playlist.play_playlist(self.playlist.current_index)
                self._manually_stopped = False
        elif event.button.id == "stop_button":
            self.playlist.player.stop()
            self._manually_stopped = True
        elif event.button.id == "prev_button":
            self.playlist.prev_song()
            self._manually_stopped = False
        elif event.button.id == "next_button":
            self.playlist.next_song()
            self._manually_stopped = False
        elif event.button.id == "random_button":
            self.playlist.random_mode = not self.playlist.random_mode
            self.playlist.upcoming_index = None
            event.button.variant = "primary" if self.playlist.random_mode else "default"
        elif event.button.id == "reroll_button":
            self.playlist.reroll_next_song()

        self._update_ui()

    def on_key(self, event) -> None:
        """Handle key presses for controlling playback."""
        if event.key == "space":
            self.on_button_pressed(
                Button.Pressed(self.query_one("#play_pause_button", Button))
            )
        elif event.key == "s":
            self.on_button_pressed(
                Button.Pressed(self.query_one("#stop_button", Button))
            )
        elif event.key == "p":
            self.on_button_pressed(
                Button.Pressed(self.query_one("#prev_button", Button))
            )
        elif event.key == "n":
            self.on_button_pressed(
                Button.Pressed(self.query_one("#next_button", Button))
            )
        elif event.key == "r":
            self.on_button_pressed(
                Button.Pressed(self.query_one("#random_button", Button))
            )
