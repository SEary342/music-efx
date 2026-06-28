from src.data import MusicEFXData
from src.ui import MusicEFX

if __name__ == "__main__":
    root_music_folder = MusicEFXData.get_default_root()
    app = MusicEFX(root_music_folder)
    app.run()
