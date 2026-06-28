import random
from pathlib import Path
from typing import Optional, TypedDict

from just_playback import Playback


class RandomMode(TypedDict):
    directory: Path
    include_nested: bool


class Player:
    def __init__(self):
        self._playback = Playback()
        self.current_song: Optional[Path] = None

    def play_song(self):
        if self._playback.active and self._playback.paused:
            self._playback.resume()
        else:
            self._playback.play()

    def load_song(self, path: Path):
        self.current_song = path
        self._playback.load_file(str(path))

    def get_progress(self) -> float:
        return self._playback.curr_pos if self._playback.active else 0.0

    def get_duration(self) -> float:
        return self._playback.duration

    def toggle_pause(self):
        if self._playback.active:
            if self._playback.paused:
                self._playback.resume()
            else:
                self._playback.pause()

    def stop(self):
        self._playback.stop()

    @property
    def is_playing(self) -> bool:
        return self._playback.playing

    @property
    def is_paused(self) -> bool:
        return self._playback.paused

    @property
    def active(self) -> bool:
        return self._playback.active


class Playlist:
    def __init__(
        self,
        playlist_file: Optional[Path] = None,
        random_init: Optional[RandomMode] = None,
    ):
        self.playlist: list[Path] = []
        self.current_index = 0
        self.player = Player()
        self.played_items: list[int] = []  # History of played indices
        self.random_mode: bool = False
        self.playlist_file = playlist_file
        self.upcoming_index: Optional[int] = None

        if playlist_file:
            self.load_playlist(playlist_file)
        elif random_init:
            self._load_from_directory(
                random_init["directory"], random_init["include_nested"]
            )
            self.random_mode = True

    def load_playlist(self, path: Path):
        # Populate the playlist list from a plain text file
        self.playlist.clear()
        self.played_items.clear()
        self.current_index = 0

        if path.exists() and path.is_file():
            with open(path, "r", encoding="utf-8") as f:
                for line in f:
                    p = Path(line.strip())
                    if (
                        p.exists()
                        and p.is_file()
                        and p.suffix.lower() in [".mp3", ".wav", ".ogg"]
                    ):
                        self.playlist.append(p)

    def _load_from_directory(self, directory: Path, include_nested: bool):
        self.playlist.clear()
        self.played_items.clear()
        self.current_index = 0

        if not directory.exists() or not directory.is_dir():
            return

        pattern = "**/*" if include_nested else "*"
        for f in sorted(directory.glob(pattern)):
            if f.is_file() and f.suffix.lower() in [".mp3", ".wav", ".ogg"]:
                self.playlist.append(f)

    def set_playlist(self, tracks: list[Path]):
        self.playlist = tracks
        self.played_items.clear()
        self.current_index = 0
        self.upcoming_index = None

    def play_playlist(self, start_index: int = 0):
        # use the player class to play audio
        if not self.playlist:
            return

        self.current_index = start_index % len(self.playlist)
        self._play_current()

    def _play_current(self):
        if (
            not self.playlist
            or self.current_index < 0
            or self.current_index >= len(self.playlist)
        ):
            return

        song = self.playlist[self.current_index]
        self.player.load_song(song)
        self.player.play_song()

        if not self.played_items or self.played_items[-1] != self.current_index:
            self.played_items.append(self.current_index)

        self.upcoming_index = None

    def _generate_upcoming(self):
        if not self.playlist:
            self.upcoming_index = None
            return

        # Check if we are going forward in history
        if self.played_items and self.current_index in self.played_items:
            idx_in_hist = self.played_items.index(self.current_index)
            if idx_in_hist < len(self.played_items) - 1:
                self.upcoming_index = self.played_items[idx_in_hist + 1]
                return

        # Otherwise, pick next song
        if self.random_mode:
            unplayed = [
                i for i in range(len(self.playlist)) if i not in self.played_items
            ]
            if unplayed:
                self.upcoming_index = random.choice(unplayed)
            else:
                self.upcoming_index = random.randint(0, len(self.playlist) - 1)
        else:
            self.upcoming_index = (self.current_index + 1) % len(self.playlist)

    def get_next_song(self) -> Optional[Path]:
        if not self.playlist:
            return None
        if self.upcoming_index is None:
            self._generate_upcoming()
        if self.upcoming_index is not None and 0 <= self.upcoming_index < len(
            self.playlist
        ):
            return self.playlist[self.upcoming_index]
        return None

    def reroll_next_song(self):
        if self.random_mode:
            self.upcoming_index = None
            self._generate_upcoming()

    def next_song(self):
        if not self.playlist:
            return

        if self.upcoming_index is None:
            self._generate_upcoming()

        if self.upcoming_index is not None:
            self.current_index = self.upcoming_index
            self._play_current()

    def prev_song(self):
        if not self.playlist or len(self.played_items) <= 1:
            return

        try:
            idx_in_hist = self.played_items.index(self.current_index)
            if idx_in_hist > 0:
                self.current_index = self.played_items[idx_in_hist - 1]
                self._play_current()
        except ValueError:
            pass
