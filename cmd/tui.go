package main

import (
	"fmt"
	"music-efx/internal/files"
	"music-efx/internal/metadata"
	"music-efx/internal/player"
	metaModel "music-efx/pkg/model"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	items         []metaModel.MP3Metadata
	selectedIndex int           // Track the currently selected file
	isPlaying     bool          // Flag to track if playback is ongoing
	currentTrack  *player.Track // Store the currently playing track
	currentFile   string        // Store the currently playing file name
	startIndex    int           // To track the pagination (first item to show)
	player        *player.Player
}

func (m *model) Init() tea.Cmd {
	// Initialize the player
	m.player = &player.Player{}
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard inputs
		switch msg.String() {
		case "up":
			// Move selection up
			if m.selectedIndex > 0 {
				m.selectedIndex--
			} else if m.startIndex > 0 {
				// If we are at the top of the current page, move the page up
				m.startIndex--
			}
		case "down":
			// Move selection down
			if m.selectedIndex < len(m.items)-1 && m.selectedIndex < m.startIndex+10 {
				m.selectedIndex++
			} else if m.selectedIndex < len(m.items)-1 {
				// If we are at the bottom of the current page, move the page down
				m.startIndex++
			}
		case "enter":
			// Select the current file to play
			selectedFile := m.items[m.selectedIndex]

			// Stop current playback if already playing
			if m.isPlaying {
				fmt.Println("Stopping current playback...")
				m.isPlaying = false
				m.currentTrack.Close() // Close the current track
				m.player.Stop()        // Stop the player
			}

			// Load the new track
			track, err := player.LoadTrack(selectedFile.Path)
			if err != nil {
				fmt.Println("Error loading track:", err)
				break
			}

			// Start playback in a new goroutine
			m.isPlaying = true
			m.currentTrack = track
			m.player.PlayTrack(track)
			m.currentTrack = track
			m.currentFile = selectedFile.Name
			fmt.Printf("Playing: %s\n", m.currentFile)

		case "esc":
			// Exit the program if no song is playing
			if !m.isPlaying {
				return m, tea.Quit
			}
			// Stop playback if a song is playing
			if m.isPlaying {
				fmt.Println("Stopping playback...")
				m.isPlaying = false
				m.currentTrack.Close()
				m.player.Stop()
			}
		case "q":
			// Stop playback and exit the program
			if m.isPlaying {
				fmt.Println("Stopping playback and quitting...")
				m.isPlaying = false
				m.currentTrack.Close()
				m.player.Stop()
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() string {
	view := "MP3 Files:\n\n"

	// Show items in pages
	for i := m.startIndex; i < len(m.items) && i < m.startIndex+10; i++ {
		item := m.items[i]
		if i == m.selectedIndex {
			// Add a clear indicator for the selected item (e.g., ">")
			view += fmt.Sprintf("> %s (%s)\n", item.Name, item.Length)
		} else {
			view += fmt.Sprintf("  %s (%s)\n", item.Name, item.Length)
		}
	}

	// Show navigation instructions
	view += "\nUse arrow keys to navigate, 'Enter' to play, 'Esc' to exit (if no song is playing), 'Q' to stop and exit.\n"

	// Only show current playback status if a file is playing
	if m.isPlaying {
		view += fmt.Sprintf("\nCurrently playing: %s\n", m.currentFile)
	}

	return view
}

func getDirectory() string {
	// Check if a directory argument is provided
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	// If no directory is provided, prompt the user for a directory
	fmt.Println("Enter directory path to scan for MP3 files:")
	var directory string
	_, err := fmt.Scanln(&directory)
	if err != nil || directory == "" {
		fmt.Println("Invalid directory path.")
		os.Exit(1)
	}

	return directory
}

func main() {
	// Get the directory path (either from arguments or prompt)
	directory := getDirectory()

	// Discover MP3 files in the specified directory
	paths, err := files.FindMP3Files(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract metadata for the MP3 files
	var metadataList []metaModel.MP3Metadata
	for _, path := range paths {
		meta, err := metadata.ExtractMetadata(path)
		if err == nil {
			metadataList = append(metadataList, meta)
		}
	}

	// Run TUI program
	if _, err := tea.NewProgram(&model{items: metadataList}).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
