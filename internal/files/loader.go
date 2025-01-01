package files

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hajimehoshi/go-mp3"
)

func LoadFile(filePath string) (*mp3.Decoder, error) {
	// Read the file content into a byte slice
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Return the error with additional context
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Convert the byte slice into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode the MP3 file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		// Return the error with additional context instead of panicking
		return nil, fmt.Errorf("mp3.NewDecoder failed for file %s: %w", filePath, err)
	}
	return decodedMp3, nil
}
