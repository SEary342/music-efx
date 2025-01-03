package player

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
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

	speaker.Init(track.Format.SampleRate, track.Format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: track.Stream, Paused: false}
	go func() {
		speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
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
		fmt.Println("No track is currently playing.")
		return
	}

	p.stopping = true
	speaker.Clear()
	p.track.Stream.Seek(0)
	p.playing = false
}

// Crossfade transitions smoothly to a new track.
func (p *Player) Crossfade(nextTrack *Track, duration time.Duration) {
	if p.track == nil || nextTrack == nil {
		fmt.Println("Both current and next tracks are required for crossfade.")
		return
	}

	fadeOut := effects.Volume{Streamer: p.track.Stream, Base: 2, Volume: -1}
	fadeIn := effects.Volume{Streamer: nextTrack.Stream, Base: 2, Volume: -1}

	combined := beep.Mix(&fadeOut, &fadeIn)

	speaker.Init(nextTrack.Format.SampleRate, nextTrack.Format.SampleRate.N(time.Second/10))
	go func() {
		speaker.Play(combined)
	}()
}
