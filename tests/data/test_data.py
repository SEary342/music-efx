from pathlib import Path
from unittest.mock import MagicMock, patch

from music_efx.data import MusicEFXData

# -------------------------
# get_default_root
# -------------------------


def test_get_default_root():
    root = MusicEFXData.get_default_root()

    assert isinstance(root, Path)
    assert root.name == "Music"


# -------------------------
# format_song_name (no tag)
# -------------------------


def test_format_song_name_strips_prefix():
    fake_path = Path("123_-_MySong.mp3")

    with patch("tinytag.TinyTag.get", side_effect=Exception("no tag")):
        result = MusicEFXData.format_song_name(fake_path)

    assert result.startswith("MySong")
    assert "123" not in result


# -------------------------
# format_song_name (with duration)
# -------------------------


def test_format_song_name_with_duration():
    fake_path = Path("Song.mp3")

    mock_tag = MagicMock()
    mock_tag.duration = 125  # 2:05

    with patch("tinytag.TinyTag.get", return_value=mock_tag):
        result = MusicEFXData.format_song_name(fake_path)

    assert "Song" in result
    assert "[2:05]" in result


# -------------------------
# get_children
# -------------------------


def test_get_children(tmp_path):
    # setup fake folder
    (tmp_path / "dirA").mkdir()
    (tmp_path / "dirB").mkdir()
    (tmp_path / "song2.mp3").write_text("x")
    (tmp_path / "song1.mp3").write_text("x")
    (tmp_path / "ignore.txt").write_text("x")

    # avoid TinyTag entirely in this test
    with patch("music_efx.data.MusicEFXData.format_song_name", side_effect=lambda p: p.stem):
        result = list(MusicEFXData(tmp_path).get_children(str(tmp_path)))

    # dirs first
    assert result[0][2] is True
    assert result[1][2] is True

    # files after
    files = result[2:]
    assert all(not f[2] for f in files)

    # sorted
    names = [f[0] for f in files]
    assert names == sorted(names, key=str.lower)
