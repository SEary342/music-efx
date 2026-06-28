from textual.widgets import Tree
from textual.widgets.tree import TreeNode

from music_efx.data import MusicEFXData


class FileSystemTree(Tree):
    """A self-contained widget that lazily loads file system data."""

    def __init__(self, data_provider: MusicEFXData, *args, **kwargs):
        # Initialize the base Tree with the root label and data
        root_label, root_path = data_provider.get_root()
        super().__init__(*args, data=root_path, **kwargs)

        self.data_provider = data_provider

    def on_mount(self) -> None:
        """Load the initial children when the widget is rendered."""
        self._load_folder_contents(self.root)
        self.root.expand()

    def on_tree_node_expanded(self, event: Tree.NodeExpanded) -> None:
        """This event is now handled inside the widget, not the main App."""
        # Stop the event from bubbling up to the App if we handle it here
        event.stop()

        if not event.node.children:
            self._load_folder_contents(event.node)

    def _load_folder_contents(self, node: TreeNode[str]) -> None:
        """Fetches data from the provider and populates the UI node."""
        if node.data is None:
            return
        for label, path, is_dir in self.data_provider.get_children(node.data):
            if is_dir:
                node.add(label, data=path)
            else:
                node.add_leaf(label, data=path)
