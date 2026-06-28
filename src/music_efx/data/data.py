import re
from collections.abc import Iterator
from pathlib import Path

from tinytag import TinyTag

AUDIO_FORMATS = {".mp3", ".wav", ".ogg"}


class MusicEFXData:
    def __init__(self, root_folder: Path):
        self.root_folder = root_folder

    @staticmethod
    def get_default_root() -> Path:
        return Path.home() / "Music"

    def get_root(self) -> tuple[str, str]:
        """Returns the label and path for the root folder."""
        return self.root_folder.name, str(self.root_folder)

    @staticmethod
    def format_song_name(item: Path) -> str:
        """Format the song path into a neat display label."""
        original_label = item.stem
        label = re.sub(r"^[^a-zA-Z]+", "", original_label)
        if not label:
            label = original_label

        try:
            tag = TinyTag.get(str(item))
            if tag.duration:
                minutes, seconds = divmod(int(tag.duration), 60)
                label = f"{label} [{minutes}:{seconds:02}]"
        except Exception:
            pass
        return label

    def get_children(self, folder_path: str) -> Iterator[tuple[str, str, bool]]:
        """
        Yields (label, path, is_dir) for a given folder.
        This is a pure Python generator—no UI dependencies.
        """
        folder = Path(folder_path)
        if not folder.is_dir():
            return

        items = []
        for item in folder.iterdir():
            if item.is_dir():
                items.append((item.name, str(item), True))
            elif item.suffix.lower() in AUDIO_FORMATS:
                label = MusicEFXData.format_song_name(item)
                items.append((label, str(item), False))

        # Sort alphabetically (case-insensitive), keeping directories at the top
        items.sort(key=lambda x: (not x[2], x[0].lower()))

        for item in items:
            yield item

    def list_audio_files_in_folder(self, folder: Path) -> list[Path]:
        """Return sorted audio files in a folder."""
        if not folder.is_dir():
            return []

        return sorted(
            [
                item
                for item in folder.iterdir()
                if item.is_file() and item.suffix.lower() in AUDIO_FORMATS
            ],
            key=lambda p: p.name.lower(),
        )

    def build_playlist_from_folder(self, folder: Path) -> list[Path]:
        return self.list_audio_files_in_folder(folder)

    def is_directory(self, path: Path) -> bool:
        return path.is_dir()
