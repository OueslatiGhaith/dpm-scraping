package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func writeResults(data ListeAttente, name string) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("Failed to marshal data: %s", err)
	}

	filename := filepath.Join(RESULTS_DIR, fmt.Sprintf("%s.json", name))
	if err := os.MkdirAll(RESULTS_DIR, 0755); err != nil {
		return fmt.Errorf("Failed to create results dir: %s", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("Failed to write file: %s", err)
	}

	log.Infof("Wrote results to %s", filename)

	return nil
}
