package src

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

func loadCheckpoint(flags *Flags) (*Checkpoint, error) {
	log.Debug("Loading checkpoint file")
	if _, err := os.Stat(checkpointFileName(flags)); os.IsNotExist(err) {
		log.Debug("Checkpoint file does not exist, creating a new one")
		return &Checkpoint{
			LastUpdated:            time.Now(),
			ProcessedGovs:          make(map[string]bool),
			ProcessedDels:          make(map[string][]string),
			PartialResultsAttente:  make(ListeAttente),
			PartialResultsOfficine: make(ListeOfficine),
		}, nil
	}

	log.Debug("Checkpoint file exists, loading it")
	data, err := os.ReadFile(checkpointFileName(flags))
	if err != nil {
		return nil, err
	}

	log.Debug("Checkpoint file loaded, unmarshalling it")
	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, err
	}

	return &checkpoint, nil
}

func saveCheckpoint(checkpoint *Checkpoint, flags *Flags) error {
	checkpoint.LastUpdated = time.Now()

	data, err := json.Marshal(checkpoint)
	if err != nil {
		return err
	}

	// create dir if needed
	dir := filepath.Dir(checkpointFileName(flags))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(checkpointFileName(flags), data, 0644)
}
