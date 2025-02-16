# Variables
PYINSTALLER_TOOL := pyinstaller
MAIN_SCRIPT := main.py

# Default target
all: uv-build pyinstaller

# Install pyinstaller as a UV tool
uv-install:
	uv tool install $(PYINSTALLER_TOOL)

# Run uv build
uv-build: uv-install
	uv build

# Run pyinstaller on the main script
pyinstaller: uv-build
	uv run $(PYINSTALLER_TOOL) --onefile $(MAIN_SCRIPT) --name music-efx

# Clean build artifacts
clean:
	rm -rf dist build *.spec
	uv uninstall $(PYINSTALLER_TOOL)

.PHONY: all uv-install uv-build pyinstaller clean
