package player

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Track struct {
	Path   string
	Length time.Duration
	Stream beep.StreamSeekCloser
	Format beep.Format
}

// LoadTrack loads an MP3 file, extracts its length, and prepares it for playback.
func LoadTrack(path string) (*Track, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	stream, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to decode mp3: %w", err)
	}

	length := time.Duration(float64(stream.Len()) / float64(format.SampleRate) * float64(time.Second))
	return &Track{Path: path, Length: length, Stream: stream, Format: format}, nil
}

// Close releases resources associated with a track.
func (t *Track) Close() {
	t.Stream.Close()
}

type Player struct {
	mu       sync.Mutex
	playing  bool
	stopping bool
	track    *Track
}

// PlayTrack starts playing a track in a non-blocking manner.
// PlayTrack starts playing a track in a non-blocking manner and displays the playback position as a progress bar.
func (p *Player) PlayTrack(track *Track) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		fmt.Println("Already playing a track. Stop it first.")
		return
	}

	p.track = track
	p.stopping = false
	p.playing = true

	// Initialize the speaker with the track's format
	speaker.Init(track.Format.SampleRate, track.Format.SampleRate.N(time.Second/10))

	// Create a control streamer to manage playback
	ctrl := &beep.Ctrl{Streamer: track.Stream, Paused: false}

	// Start a goroutine to display the playback position and progress bar
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond) // Update every 0.5 seconds
		defer ticker.Stop()

		for range ticker.C {
			p.mu.Lock()
			if !p.playing || p.stopping || p.track == nil {
				p.mu.Unlock()
				return
			}

			// Calculate the current position and total length in seconds
			position := float64(p.track.Stream.Position()) / float64(p.track.Format.SampleRate)
			total := float64(p.track.Stream.Len()) / float64(p.track.Format.SampleRate)
			p.mu.Unlock()

			// Calculate the progress percentage
			progress := position / total

			// Construct the progress bar
			barLength := 50 // Length of the progress bar (in characters)
			progressBar := ""
			for i := 0; i < barLength; i++ {
				if float64(i)/float64(barLength) < progress {
					progressBar += "="
				} else {
					progressBar += " "
				}
			}

			// Clear the line and update the playback position with progress bar
			fmt.Printf("\rPlaying: %s [%s] %.1fs / %.1fs", p.track.Path, progressBar, position, total)
			// Flush the output to ensure it prints immediately
			fmt.Print()
		}
	}()

	// Start playing the track
	go func() {
		speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
			// Callback when the track ends
			p.mu.Lock()
			p.playing = false
			p.mu.Unlock()
		})))
	}()
}

// Stop stops the currently playing track.
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return
	}

	p.stopping = true
	speaker.Clear()
	p.track.Stream.Seek(0)
	p.playing = false
}
