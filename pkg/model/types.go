package model

import "time"

type MP3Metadata struct {
	Name   string
	Length time.Duration
	Path   string
}
