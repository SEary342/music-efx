package model

import (
	"time"

	"gopkg.in/yaml.v3"
)

type MP3Metadata struct {
	Name   string
	Length time.Duration
	Path   string
}

type PlaylistData struct {
	Name string `yaml:"Name"`
	End  Time   `yaml:"End"`
	Path string `yaml:"Path"`
}

// Time wraps time.Time to implement custom unmarshaling.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalYAML(value *yaml.Node) error {
	// Parse the time string in "HH:mm" format
	parsedTime, err := time.Parse("15:04", value.Value)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}
