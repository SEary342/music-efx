package player

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
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
	playing bool
	ctrl    *beep.Ctrl
	track   *Track
}

func (p *Player) PlayTrack(track *Track) {
	speaker.Lock()
	defer speaker.Unlock()

	if p.playing {
		fmt.Println("Already playing a track. Stop it first.")
		return
	}

	p.track = track
	p.playing = true

	// Initialize the speaker with the track's format
	speaker.Init(track.Format.SampleRate, track.Format.SampleRate.N(time.Second/10))

	// Create a control streamer to manage playback
	p.ctrl = &beep.Ctrl{Streamer: track.Stream, Paused: false}

	go func() {
		speaker.Play(beep.Seq(p.ctrl, beep.Callback(func() {
			// Callback when the track ends
			speaker.Lock()
			p.playing = false
			speaker.Unlock()
		})))
	}()
}

// Stop stops the currently playing track.
func (p *Player) Stop() {
	speaker.Lock()
	defer speaker.Unlock()

	if !p.playing {
		return
	}
	p.ctrl.Paused = true
	p.playing = false
	p.track.Close()
}
