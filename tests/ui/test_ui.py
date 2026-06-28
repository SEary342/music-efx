from pathlib import Path
from unittest.mock import MagicMock, patch

# ----------------------------
# helper fixture
# ----------------------------


def make_app(app):
    app.query_one = MagicMock()
    app.query = MagicMock()
    return app


# ----------------------------
# tree selection: file
# ----------------------------


def test_file_selection_sets_playlist_and_plays(app):
    app = make_app(app)

    node = MagicMock()
    node.data = "/music/album/a.mp3"

    event = MagicMock()
    event.node = node

    fake_files = [
        Path("/music/album/a.mp3"),
        Path("/music/album/b.mp3"),
    ]

    fake_parent = MagicMock()
    fake_parent.glob.return_value = fake_files

    fake_path = MagicMock()
    fake_path.suffix.lower.return_value = ".mp3"
    fake_path.parent = fake_parent

    with (
        patch("pathlib.Path", return_value=fake_path),
        patch.object(app.playlist, "set_playlist") as set_playlist,
        patch.object(app.playlist, "play_playlist") as play,
    ):
        app.on_tree_node_selected(event)

        set_playlist.assert_called_once()
        play.assert_called_once()
        assert app._manually_stopped is False


# ----------------------------
# play/pause button
# ----------------------------


def test_play_pause_button_starts_playback(app):
    app = make_app(app)

    app.playlist.player = MagicMock()
    app.playlist.player.is_playing = False
    app.playlist.player.is_paused = False
    app.playlist.player.current_song = None
    app.playlist.player.active = False

    # IMPORTANT: prevent MagicMock math explosion
    app.playlist.player.get_duration.return_value = 0
    app.playlist.player.get_progress.return_value = 0

    btn = MagicMock()
    btn.id = "play_pause_button"

    event = MagicMock(button=btn)

    with patch.object(app.playlist, "play_playlist") as play:
        app.on_button_pressed(event)
        play.assert_called_once()


# ----------------------------
# stop button
# ----------------------------


def test_stop_button_stops_player(app):
    app = make_app(app)

    btn = MagicMock()
    btn.id = "stop_button"

    event = MagicMock(button=btn)

    with patch.object(app.playlist.player, "stop") as stop:
        app.on_button_pressed(event)
        stop.assert_called_once()
        assert app._manually_stopped is True


# ----------------------------
# next button
# ----------------------------


def test_next_button_advances(app):
    app = make_app(app)

    btn = MagicMock()
    btn.id = "next_button"

    event = MagicMock(button=btn)

    with patch.object(app.playlist, "next_song") as next_song:
        app.on_button_pressed(event)
        next_song.assert_called_once()


# ----------------------------
# random toggle
# ----------------------------


def test_random_toggle(app):
    app = make_app(app)

    btn = MagicMock()
    btn.id = "random_button"
    btn.variant = "default"

    event = MagicMock(button=btn)

    app.playlist.random_mode = False

    app.on_button_pressed(event)

    assert app.playlist.random_mode is True
