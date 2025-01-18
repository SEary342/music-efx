package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"music-efx/internal/config"
	"music-efx/internal/files"
	"music-efx/internal/metadata"
	"music-efx/internal/player"
	"music-efx/internal/playlist"
	metaModel "music-efx/pkg/model"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// Menu-related fields
	menuItems     []string
	menuIndex     int
	inPlaylist    bool
	inFolderNav   bool
	inSearch      bool
	searchQuery   string
	searchResults []metaModel.MP3Metadata
	playlists     []metaModel.PlaylistData
	// Folder navigation-related fields
	directoryTree map[string][]metaModel.MP3Metadata
	currentDir    string
	// MP3 selection-related fields
	items         []metaModel.MP3Metadata
	allItems      []metaModel.MP3Metadata
	selectedIndex int
	// Playback-related fields
	isPlaying   bool
	currentFile string
	player      *player.Player
	// Stop channel for playback control
	stopChan chan bool
}

func (m *model) Init() tea.Cmd {
	// Initialize the player and stop channel
	m.player = &player.Player{}
	m.stopChan = make(chan bool) // Initialize the stop channel
	// Load main menu items
	m.reset()
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key input
		switch msg.String() {
		case "up":
			// Handle up arrow key press
			if m.inSearch || m.inPlaylist || m.inFolderNav {
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			} else {
				if m.menuIndex > 0 {
					m.menuIndex--
				}
			}
		case "down":
			// Handle down arrow key press
			if m.inSearch || m.inPlaylist || m.inFolderNav {
				if m.selectedIndex < len(m.playlists)-1 {
					m.selectedIndex++
				}
			} else {
				if m.menuIndex < len(m.menuItems)-1 {
					m.menuIndex++
				}
			}
		case "enter":
			if m.inPlaylist {
				// Handle playlist selection
				playlistMeta := m.playlists[m.selectedIndex]
				mp3Meta, err := metadata.LoadMp3Metadata(playlistMeta.Path)
				if err != nil {
					fmt.Println("Failed to load playlist mp3 files.")
					return m, nil
				}
				go playlist.RandomPlay(mp3Meta, m.handlePlayback)
			} else if m.inFolderNav {
				// Handle folder navigation
				if m.selectedIndex < len(m.items) {
					selectedFile := m.items[m.selectedIndex]
					if selectedFile.Name == selectedFile.Path {
						m.currentDir = selectedFile.Path
						m.updateItemsForCurrentDir()
					} else {
						m.handlePlayback(selectedFile)
					}
				}
			} else if m.inSearch {
				// Handle search result selection
				if m.selectedIndex < len(m.searchResults) {
					selectedFile := m.searchResults[m.selectedIndex]
					m.handlePlayback(selectedFile)
				}
			} else {
				// Handle main menu actions
				switch m.menuItems[m.menuIndex] {
				case "Auto-Playlist":
					go m.startAutoPlaylist()
				case "Playlists":
					m.inPlaylist = true
					m.menuItems = []string{"Back"}
					for _, playlist := range m.playlists {
						m.menuItems = append(m.menuItems, playlist.Name)
					}
				case "Folder Navigation":
					m.inFolderNav = true
					m.updateItemsForCurrentDir()
				case "Search":
					m.inSearch = true
					m.searchQuery = ""
					m.searchResults = nil
				case "Quit":
					return m, tea.Quit
				}
			}
		case "esc":
			// Handle stopping the auto-playlist with ESC
			if m.isPlaying {
				fmt.Println("Stopping playback...")
				m.isPlaying = false
				m.player.Stop()
				fmt.Print("\033[H\033[2J") // Clear terminal screen
				m.reset()
			} else if m.inPlaylist || m.inFolderNav || m.inSearch {
				// If we are in a playlist, folder navigation, or search, reset
				fmt.Print("\033[H\033[2J") // Clear terminal screen
				m.reset()
			} else {
				// Exit the program if in the main menu
				return m, tea.Quit
			}

		case "backspace":
			// Handle backspace for search query
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
			m.updateSearchResults()
		case "space":
			// Handle space character for search query
			m.searchQuery += " "
			m.updateSearchResults()

		default:
			// Handle typing (any other keys)
			if msg.String() != "" && len(msg.String()) == 1 {
				m.searchQuery += msg.String()
				m.updateSearchResults()
			}
		}
	}

	return m, nil
}

func (m *model) startAutoPlaylist() {
	// Run the playlist logic in a goroutine to avoid blocking the UI
	go func() {
		playlistMeta := make([]metaModel.PlaylistData, len(m.playlists))
		copy(playlistMeta, m.playlists)

		sort.Slice(playlistMeta, func(i, j int) bool {
			return playlistMeta[i].End.Before(playlistMeta[j].End.Time)
		})

		mp3MetaMap := make(map[string][]metaModel.MP3Metadata)
		for _, lst := range playlistMeta {
			mp3Meta, err := metadata.LoadMp3Metadata(lst.Path)
			if err != nil {
				fmt.Println("Failed to load playlist mp3 files.")
				continue
			}
			mp3MetaMap[lst.Name] = mp3Meta
		}

		for _, lst := range playlistMeta {
			duration := int(time.Until(lst.End.Time).Seconds())
			if duration <= 0 {
				fmt.Println("Skipping expired playlist:", lst.Name)
				continue
			}

			// Stop the currently playing track, if any
			m.player.Stop()

			// Start the next playlist
			fmt.Println("Starting playlist:", lst.Name)
			go playlist.GenerateAndPlay(mp3MetaMap[lst.Name], duration, m.handlePlayback)

			// Wait for the playlist to finish or until stop signal is received
			time.Sleep(time.Duration(duration) * time.Second)
		}

		m.player.Stop()
		fmt.Println("All playlists have finished.")
	}()
}

func (m *model) reset() {
	m.inPlaylist = false
	m.inFolderNav = false
	m.inSearch = false
	m.menuItems = []string{"Auto-Playlist", "Playlists", "Folder Navigation", "Search", "Quit"}
}

func (m *model) handlePlayback(file metaModel.MP3Metadata) {
	if m.isPlaying {
		fmt.Println("Stopping current playback...")
		m.isPlaying = false
		m.player.Stop()
	}

	// Clear terminal screen after playback
	fmt.Print("\033[H\033[2J")

	// Load the new track
	track, err := player.LoadTrack(file.Path)
	if err != nil {
		fmt.Println("Error loading track:", err)
		return
	}

	// Start playback in a new goroutine
	m.isPlaying = true
	m.player.PlayTrack(track)
	m.currentFile = file.Name
}

func (m *model) updateSearchResults() {
	// Find files that match the search query (simple contains check)
	m.searchResults = nil
	//fmt.Println(m.allItems)
	for _, file := range m.allItems {
		//fmt.Println(m.searchQuery)
		if strings.Contains(strings.ToLower(file.Name), strings.ToLower(m.searchQuery)) {
			m.searchResults = append(m.searchResults, file)
		}
	}
}

func (m *model) updateItemsForCurrentDir() {
	// Clear current items
	m.items = nil

	// Get the MP3 files for the current directory
	filesInDir, exists := m.directoryTree[m.currentDir]
	if exists {
		// Add the files in the current directory to the items list
		m.items = filesInDir
		// Add subdirectories as folder navigation options
		for dir := range m.directoryTree {
			if strings.HasPrefix(dir, m.currentDir+"/") && dir != m.currentDir {
				m.items = append(m.items, metaModel.MP3Metadata{Path: dir, Name: dir})
			}
		}
	} else {
		// If no files, show a "back" option
		m.items = append(m.items, metaModel.MP3Metadata{Path: "Back", Name: "Go Up"})
	}
}

func (m *model) View() string {
	var view string

	if m.inSearch {
		view = "Search for MP3:\n" + m.searchQuery + "\n"
		view += "Results:\n"
		view += m.renderMenu(m.searchResults, m.selectedIndex)
	} else if m.inPlaylist {
		view = "Select a Playlist:\n"
		view += m.renderMenu(m.playlists, m.selectedIndex)
	} else if m.inFolderNav {
		view = "Select a Folder/File:\n"
		view += m.renderMenu(m.items, m.selectedIndex)
	} else {
		view = "Main Menu:\n"
		view += m.renderMenu(m.menuItems, m.menuIndex)
	}

	return view
}

func (m *model) renderMenu(items interface{}, selectedIndex int) string {
	var view string

	switch items := items.(type) {
	case []metaModel.MP3Metadata:
		// Render MP3 items
		for i, item := range items {
			dispName := item.Name // Add folder icon for directories
			if item.Name == item.Path {
				nameParts := strings.Split(item.Name, "/")
				dispName = "📁 " + nameParts[len(nameParts)-1] // Add folder icon for directories
			}
			if i == selectedIndex {
				view += "> " + dispName + "\n" // Add ">" for the selected item
			} else {
				view += "  " + dispName + "\n"
			}
		}
	case []metaModel.PlaylistData:
		// Render MP3 items
		for i, item := range items {
			if i == selectedIndex {
				view += "> " + item.Name + "\n" // Add ">" for the selected item
			} else {
				view += "  " + item.Name + "\n"
			}
		}
	case []string:
		// Render menu items (directories and files)
		for i, item := range items {
			if i == selectedIndex {
				view += "> " + item + "\n" // Add ">" for the selected item
			} else {
				view += "  " + item + "\n"
			}
		}
	}

	return view
}

func main() {
	// Get the directory path (either from arguments or prompt)
	directory := files.GetDirectory()

	metadataList, err := metadata.LoadMp3Metadata(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a directory tree from the file paths
	directoryTree := files.CreateDirectoryTree(metadataList)

	playlists := config.LoadPlaylistYaml()

	// Run TUI program
	if _, err := tea.NewProgram(&model{
		directoryTree: directoryTree,
		currentDir:    directory,
		allItems:      metadataList,
		playlists:     playlists,
	}).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
