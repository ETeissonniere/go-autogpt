package agent

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eteissonniere/hercules/llms"
)

type Exporter interface {
	Export(llms.ChatMessage) error
}

type doNotExport struct{}

func (e *doNotExport) Export(_ llms.ChatMessage) error {
	return nil
}

func DoNotExport() Exporter {
	return &doNotExport{}
}

type exportToFile struct {
	os.File
}

func (e *exportToFile) Export(msg llms.ChatMessage) error {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = e.Write(append(marshalled, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func ExportToFile(path string) (Exporter, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return &exportToFile{*file}, nil
}
