package config

import (
	"music-efx/pkg/model"
	"os"
	"reflect"
	"testing"
	"time"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestLoadConfig(t *testing.T) {
	// Parse times using time.Parse
	time1, _ := time.Parse("15:04", "10:00")
	time2, _ := time.Parse("15:04", "10:30")

	tests := []struct {
		name         string
		fileContent  string
		expectError  bool
		expectedData []model.PlaylistData
	}{
		{
			name: "Valid Config",
			fileContent: `
---
- Name: Segment1
  End: "10:00"
  Path: "/path/to/segment1/"
- Name: Segment2
  End: "10:30"
  Path: "/path/to/segment2/"
`,
			expectError: false,
			expectedData: []model.PlaylistData{
				// Convert time.Time to model.Time by embedding the time.Time into model.Time
				{Name: "Segment1", End: model.Time{Time: time1}, Path: "/path/to/segment1/"},
				{Name: "Segment2", End: model.Time{Time: time2}, Path: "/path/to/segment2/"},
			},
		},
		{
			name:        "File Does Not Exist",
			fileContent: "",
			expectError: true,
		},
		{
			name: "Invalid YAML",
			fileContent: `
---
- Name: Segment1
  End: "10:00"
  Path: "/path/to/segment1/"
- Name Segment2
  End: "10:30"
  Path: "/path/to/segment2/"
`,
			expectError: true,
		},
		{
			name:        "Empty File",
			fileContent: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFilePath string
			if tt.fileContent != "" {
				tmpFilePath = createTempFile(t, tt.fileContent)
				defer os.Remove(tmpFilePath) // Clean up the file after the test
			} else {
				tmpFilePath = "/non-existent-file.yaml"
			}

			gotData, err := loadConfig(tmpFilePath)
			if (err != nil) != tt.expectError {
				t.Errorf("LoadConfig() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError && !reflect.DeepEqual(gotData, &tt.expectedData) {
				t.Errorf("LoadConfig() = %v, expected %v", gotData, tt.expectedData)
			}
		})
	}
}
