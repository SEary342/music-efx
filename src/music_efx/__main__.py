from music_efx.data import MusicEFXData
from music_efx.ui import MusicEFX


def main():
    root_music_folder = MusicEFXData.get_default_root()
    app = MusicEFX(root_music_folder)
    app.run()


if __name__ == "__main__":
    main()
