package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func writeResults(data *Checkpoint, flags *Flags) error {
	var err error
	var jsonData []byte

	switch flags.Mode {
	case MODE_ATTENTE:
		jsonData, err = json.MarshalIndent(data.PartialResultsAttente, "", "\t")
	case MODE_OFFICINE:
		jsonData, err = json.MarshalIndent(data.PartialResultsOfficine, "", "\t")
	}

	if err != nil {
		return fmt.Errorf("Failed to marshal data: %s", err)
	}

	filename := filepath.Join(RESULTS_DIR, resultsFileName(flags))
	if err := os.MkdirAll(RESULTS_DIR, 0755); err != nil {
		return fmt.Errorf("Failed to create results dir: %s", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("Failed to write file: %s", err)
	}

	log.Infof("Wrote results to %s", filename)

	return nil
}
