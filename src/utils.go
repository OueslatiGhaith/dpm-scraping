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

func navigateToGouvernourat(page *rod.Page, link string, gouvernourat Gouvernourat, flags *Flags) error {
	// navigate to the main page and select the gouvernourat
	log.Debugf("Navigating to gouvernourat %s", gouvernourat)
	if err := page.Navigate(link); err != nil {
		return fmt.Errorf("failed to navigate to main page: %w", err)
	}

	// select gouvernourat and JOUR option
	log.Debugf("Selecting gouvernourat %s", gouvernourat)
	page.MustElement("select[name='cod_gouv']").MustSelect(string(gouvernourat))

	switch flags.Time {
	case TIME_JOUR:
		log.Debugf("Selecting time JOUR")
		page.MustElement("input[value='ON']").MustClick()
	case TIME_NUIT:
		log.Debugf("Selecting time NUIT")
		page.MustElement("input[value='OFF']").MustClick()
	default:
		return fmt.Errorf("invalid time: %d", flags.Time)
	}

	log.Debugf("Submitting form")
	page.MustElement("input[type='submit']").MustClick()

	// wait for the page to be ready
	log.Debugf("Waiting for page to be ready")
	page.MustWaitLoad()

	log.Debugf("Page ready")
	return nil
}

// func trimWhitespace(v interface{}) {
// 	val := reflect.ValueOf(v)

// 	if val.Kind() == reflect.Pointer {
// 		val = val.Elem()
// 	}

// 	if !val.IsValid() || !val.CanSet() {
// 		return
// 	}

// 	switch val.Kind() {
// 	case reflect.String:
// 		log.Error("Hello")
// 		s := val.String()
// 		s = strings.TrimSpace(s)
// 		s = strings.ReplaceAll(s, " ", "")
// 		if strings.Contains(s, " ") {
// 			log.Fatal("Found non-breaking space in string")
// 		}
// 		val.SetString(s)

// 	case reflect.Struct:
// 		for i := 0; i < val.NumField(); i++ {
// 			field := val.Field(i)
// 			if field.CanInterface() {
// 				trimWhitespace(field.Addr().Interface())
// 			}
// 		}

// 	case reflect.Slice, reflect.Array:
// 		for i := 0; i < val.Len(); i++ {
// 			item := val.Index(i)
// 			if item.CanAddr() && item.Addr().CanInterface() {
// 				trimWhitespace(item.Addr().Interface())
// 			}
// 		}

// 	case reflect.Map:
// 		for _, key := range val.MapKeys() {
// 			item := val.MapIndex(key)
// 			if item.CanInterface() {
// 				newVal := reflect.New(item.Type()).Elem()
// 				newVal.Set(item)
// 				if newVal.CanAddr() {
// 					trimWhitespace(newVal.Addr().Interface())
// 					val.SetMapIndex(key, newVal)
// 				}
// 			}
// 		}
// 	}
// }
