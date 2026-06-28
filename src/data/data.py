from typing import Iterator
from pathlib import Path


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

    def get_children(self, folder_path: str) -> Iterator[tuple[str, str, bool]]:
        """
        Yields (label, path, is_dir) for a given folder.
        This is a pure Python generator—no UI dependencies.
        """
        folder = Path(folder_path)
        if not folder.is_dir():
            return

        for item in sorted(folder.iterdir()):
            if item.is_dir():
                yield item.name, str(item), True
            elif item.suffix.lower() in AUDIO_FORMATS:
                yield item.name, str(item), False
