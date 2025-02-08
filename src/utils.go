package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
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

func getByXPath(page *rod.Page, xPath string) string {
	return page.MustEval(fmt.Sprintf(`function() { 
		return document.evaluate("%s", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.textContent; 
	}`, xPath)).String()
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to convert %s to int: %s", s, err)
	}
	return i
}
