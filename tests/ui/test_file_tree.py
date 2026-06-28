from unittest.mock import MagicMock, patch

import pytest

from music_efx.ui.widgets.file_tree import FileSystemTree

# ----------------------------
# shared fixture
# ----------------------------


@pytest.fixture
def provider():
    p = MagicMock()
    p.get_root.return_value = ("Music", "/music")
    return p


# ----------------------------
# on_mount
# ----------------------------


def test_on_mount_loads_root_children(provider):
    with patch("textual.widgets.Tree.__init__", return_value=None):
        tree = FileSystemTree(provider)

        tree.root = MagicMock()
        tree._load_folder_contents = MagicMock()

        tree.on_mount()

        tree._load_folder_contents.assert_called_once_with(tree.root)
        tree.root.expand.assert_called_once()


# ----------------------------
# expand loads children when empty
# ----------------------------


def test_on_node_expanded_loads_children_when_empty(provider):
    node = MagicMock()
    node.children = []

    event = MagicMock()
    event.node = node

    with patch("textual.widgets.Tree.__init__", return_value=None):
        tree = FileSystemTree(provider)

        tree._load_folder_contents = MagicMock()

        tree.on_tree_node_expanded(event)

        event.stop.assert_called_once()
        tree._load_folder_contents.assert_called_once_with(node)


# ----------------------------
# expand skips when already loaded
# ----------------------------


def test_on_node_expanded_skips_when_children_exist(provider):
    node = MagicMock()
    node.children = [MagicMock()]

    event = MagicMock()
    event.node = node

    with patch("textual.widgets.Tree.__init__", return_value=None):
        tree = FileSystemTree(provider)

        tree._load_folder_contents = MagicMock()

        tree.on_tree_node_expanded(event)

        event.stop.assert_called_once()
        tree._load_folder_contents.assert_not_called()


# ----------------------------
# loading folder contents
# ----------------------------


def test_load_folder_contents_adds_nodes(provider):
    provider.get_children.return_value = [
        ("FolderA", "/a", True),
        ("Song1", "/song.mp3", False),
    ]

    with patch("textual.widgets.Tree.__init__", return_value=None):
        tree = FileSystemTree(provider)

        node = MagicMock()
        node.data = "/root"
        node.add = MagicMock()
        node.add_leaf = MagicMock()

        tree._load_folder_contents(node)

        provider.get_children.assert_called_once_with("/root")
        node.add.assert_called_once_with("FolderA", data="/a")
        node.add_leaf.assert_called_once_with("Song1", data="/song.mp3")


# ----------------------------
# skips when node has no data
# ----------------------------


def test_load_folder_contents_skips_none_data(provider):
    with patch("textual.widgets.Tree.__init__", return_value=None):
        tree = FileSystemTree(provider)

        node = MagicMock()
        node.data = None

        tree._load_folder_contents(node)

        provider.get_children.assert_not_called()
