from pathlib import Path
import random
import pygame
from textual.app import App, ComposeResult
from textual.containers import Horizontal, Vertical
from textual.widgets import Tree, Static, Button, Footer, Header, ProgressBar
import time
import threading


class MusicEFX(App):
    """TUI Music Player Application."""

    CSS = """
        ProgressBar {
            margin-right: 2;
        }

        #folder_tree {
            height: 70%;
        }

        .pressed {
            background: blue 80%;
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
        self.root_folder = root_folder
        self.current_playlist = []
        self.current_index = 0
        self.is_playing = False
        self.is_paused = False
        self.random_mode = False
        self.played_items = dict()
        self.song_length = 0
        self.progress_thread = None
        self.start_time = 0  # Track start time for ETA calculation
        pygame.mixer.init()

    def compose(self) -> ComposeResult:
        yield Header()
        yield Tree(label="Music Folders", data=str(self.root_folder), id="folder_tree")
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
        yield Footer()

    def on_mount(self) -> None:
        # Populate the folder tree
        tree = self.query_one("#folder_tree", Tree)
        tree.root.expand()
        self.populate_tree(tree.root, self.root_folder)

    def populate_tree(self, node, folder: Path) -> None:
        """Recursively add folders and music files to the tree."""
        for item in sorted(folder.iterdir()):
            if item.is_dir():
                branch = node.add(item.name, data=str(item))
                self.populate_tree(branch, item)
            elif item.suffix.lower() in [".mp3", ".wav", ".ogg"]:
                node.add_leaf(item.name, data=str(item))

    def play_song(self):
        """Play the current song."""
        if self.current_playlist:
            song_path = self.current_playlist[self.current_index]
            self.played_items[self.current_index] = True
            self.enable_buttons()
            pygame.mixer.music.load(song_path)
            pygame.mixer.music.play()
            self.is_playing = True
            self.is_paused = False
            self.song_length = pygame.mixer.Sound(song_path).get_length()
            self.start_time = time.time()  # Track the start time
            song_status = self.query_one("#song_status", Static)
            song_status.update(f"\nPlaying: {Path(song_path).name}")
            self.reset_progress_bar()
            self.update_progress_bar()

            # Highlight the current song in the file picker
            tree = self.query_one("#folder_tree", Tree)
            for node in tree.query():
                if node.data == str(song_path):
                    node.highlight()

    def reset_progress_bar(self):
        """Reset the progress bar to its initial state."""
        progress_bar = self.query_one("#progress_bar", ProgressBar)
        progress_bar.progress = 0

    def update_progress_bar(self):
        """Update the progress bar and ETA while the song is playing."""

        def progress_worker():
            progress_bar = self.query_one("#progress_bar", ProgressBar)
            eta_display = self.query_one(
                "#eta_display", Static
            )  # Static widget to show ETA

            progress_bar.update(total=int(self.song_length))

            while self.is_playing and pygame.mixer.music.get_busy():
                elapsed = time.time() - self.start_time
                remaining_time = self.song_length - elapsed
                progress_bar.total = 100
                current_step = int((elapsed / self.song_length) * 100)
                progress_bar.progress = min(max(current_step, 0), 100)

                # Update ETA display
                minutes, seconds = divmod(int(remaining_time), 60)
                eta_display.update(f"{minutes:02}:{seconds:02}")

                time.sleep(0.5)
            if not self.is_paused and self.is_playing:
                self.call_later(self.next_song)

        if self.progress_thread and self.progress_thread.is_alive():
            return
        self.progress_thread = threading.Thread(target=progress_worker, daemon=True)
        self.progress_thread.start()

    def stop_song(self):
        """Stop the current song."""
        pygame.mixer.music.stop()
        self.is_playing = False
        self.is_paused = False
        song_status = self.query_one("#song_status", Static)
        song_status.update("\nNo song playing")
        self.reset_progress_bar()

        # Clear the ETA when the song is stopped
        eta_display = self.query_one("#eta_display", Static)
        eta_display.update("")

    def pause_or_resume_song(self):
        """Toggle between pause and resume for the current song."""
        if self.is_playing:
            if self.is_paused:
                pygame.mixer.music.unpause()
                self.is_paused = False
                song_status = self.query_one("#song_status", Static)
                song_status.update(
                    f"\nPlaying: {Path(self.current_playlist[self.current_index]).name}"
                )
            else:
                pygame.mixer.music.pause()
                self.is_paused = True
                song_status = self.query_one("#song_status", Static)
                song_status.update(
                    f"\nPaused: {Path(self.current_playlist[self.current_index]).name}"
                )

    def next_song(self):
        """Automatically move to the next song in the playlist."""
        if self.current_playlist and len(self.played_items) < len(
            self.current_playlist
        ):
            hist_list = list(self.played_items.keys())
            idx_in_hist = hist_list.index(self.current_index)
            old_items = hist_list[:-1]
            if self.current_index not in old_items:
                if self.random_mode:
                    curr_index = self.current_index
                    new_index = curr_index
                    while new_index == curr_index and new_index in self.played_items:
                        new_index = random.randint(0, len(self.current_playlist) - 1)
                    self.current_index = new_index
                else:
                    self.current_index = (self.current_index + 1) % len(
                        self.current_playlist
                    )
            else:
                self.current_index = hist_list[idx_in_hist+1]
            self.reset_progress_bar()
            self.play_song()

    def on_tree_node_selected(self, event: Tree.NodeSelected) -> None:
        """Handle folder or file selection in the tree."""
        selected_path = Path(event.node.data)
        self.enable_buttons()
        if selected_path.is_dir():
            self.current_playlist = [str(f) for f in selected_path.glob("*.mp3")]
            self.current_index = 0
            self.played_items = dict()
        elif selected_path.suffix.lower() in [".mp3", ".wav", ".ogg"]:
            # Include all audio files in the directory for playlist context
            self.current_playlist = [str(f) for f in selected_path.parent.glob("*.mp3")]
            self.current_index = self.current_playlist.index(str(selected_path))
            self.reset_progress_bar()
            self.play_song()

    def enable_buttons(self) -> None:
        if self.current_playlist:
            buttons: list[Button] = self.query(".btn")
            for btn in buttons:
                if btn.id == "prev_button" and len(self.played_items) <= 1:
                    btn.disabled = True
                else:
                    btn.disabled = False

    def on_button_pressed(self, event: Button.Pressed) -> None:
        """Handle button presses."""
        if event.button.id == "play_pause_button":
            if not self.is_playing:
                self.play_song()
            else:
                self.pause_or_resume_song()
        elif event.button.id == "stop_button":
            self.stop_song()
        elif event.button.id == "prev_button":
            if self.current_playlist and len(self.played_items.keys()) > 1:
                dict_key = list(self.played_items.keys()).index(self.current_index) - 1
                self.current_index = list(self.played_items.keys())[dict_key]
                self.reset_progress_bar()
                self.play_song()
        elif event.button.id == "next_button":
            self.next_song()  # Call the new next_song method
        elif event.button.id == "random_button":
            if self.current_playlist:
                self.random_mode = not self.random_mode
                event.button.classes = "pressed" if self.random_mode else ""

    def on_key(self, event) -> None:
        """Handle key presses for controlling playback."""
        if event.key == "space":
            self.pause_or_resume_song()
        elif event.key == "s":
            self.stop_song()
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


if __name__ == "__main__":
    root_music_folder = Path.home() / "Music"
    app = MusicEFX(root_music_folder)
    app.run()
