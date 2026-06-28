# tests/conftest.py
from pathlib import Path

import pytest

from music_efx.ui import MusicEFX


@pytest.fixture
def app():
    return MusicEFX(Path("/music"))