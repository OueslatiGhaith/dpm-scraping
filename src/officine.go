package src

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
)

func scrapeOfficines(ctx context.Context, page *rod.Page, checkpoint *Checkpoint, flags *Flags) error {
	log.Info("Navigating to officine page")
	if err := page.Navigate(OFFICINE); err != nil {
		return fmt.Errorf("failed to navigate to officine page: %w", err)
	}

	log.Debug("Getting list of governments")
	govSelect := page.MustElement("select[name='cod_gouv']")
	options := govSelect.MustElements("option")

	var gouvernourats []Gouvernourat
	for _, opt := range options {
		gov := opt.MustText()
		gouvernourats = append(gouvernourats, Gouvernourat(gov))
	}

	for _, gov := range gouvernourats {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
		}

		if _, ok := checkpoint.ProcessedGovs[gov]; ok {
			log.Warnf("\tGovernment %s already processed, skipping", gov)
			continue
		}

		log.Infof("Processing government %s", gov)
		checkpoint.CurrentGov = gov

		if err := processOfficines(ctx, page, gov, checkpoint, flags); err != nil {
			return fmt.Errorf("failed to process government %s: %w", gov, err)
		}

		checkpoint.ProcessedGovs[gov] = true
		checkpoint.CurrentGov = ""
		checkpoint.CurrentDel = ""

		// Save checkpoint after each government
		if err := saveCheckpoint(checkpoint, flags); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	return nil
}

func processOfficines(ctx context.Context, page *rod.Page, gov Gouvernourat, checkpoint *Checkpoint, flags *Flags) error {
	if err := navigateToGouvernourat(page, OFFICINE, gov, flags); err != nil {
		return fmt.Errorf("failed to navigate to main page: %w", err)
	}

	delSelect := getDelegationSelect(page, flags)
	options := delSelect.MustElements("option")

	log.Debug("Getting list of delegations")
	var delegations []Delegation
	for _, opt := range options {
		del := opt.MustText()
		delegations = append(delegations, Delegation(del))
	}

	if checkpoint.PartialResultsOfficine[gov] == nil {
		checkpoint.PartialResultsOfficine[gov] = make(map[Delegation][]*Officine)
	}

	for _, del := range delegations {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
		}

		if slices.Contains(checkpoint.ProcessedDels[gov], del) {
			log.Warnf("Delegation %s already processed, skipping", del)
			continue
		}

		log.Infof("\tProcessing delegation %s", del)
		if err := processOfficineDelegation(page, gov, del, checkpoint, flags); err != nil {
			if saveErr := saveCheckpoint(checkpoint, flags); saveErr != nil {
				log.Errorf("Failed to save checkpoint: %s", saveErr)
			}
			return fmt.Errorf("failed to process delegation %s: %w", del, err)
		}

		if checkpoint.ProcessedDels[gov] == nil {
			checkpoint.ProcessedDels[gov] = make([]Delegation, 0)
		}
		checkpoint.ProcessedDels[gov] = append(checkpoint.ProcessedDels[gov], del)

		if err := saveCheckpoint(checkpoint, flags); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	return nil
}

func processOfficineDelegation(page *rod.Page, gov Gouvernourat, del Delegation, checkpoint *Checkpoint, flags *Flags) error {
	if err := navigateToGouvernourat(page, OFFICINE, gov, flags); err != nil {
		return fmt.Errorf("failed to navigate to main page: %w", err)
	}

	page.MustWaitLoad()
	getDelegationSelect(page, flags).MustSelect(string(del))
	page.MustElement("input[type='submit']").MustClick()

	page.MustWaitLoad()

	officines, err := extractOfficines(page)
	if err != nil {
		return fmt.Errorf("failed to extract officines: %w", err)
	}

	if checkpoint.PartialResultsOfficine[gov] == nil {
		checkpoint.PartialResultsOfficine[gov] = make(map[Delegation][]*Officine)
	}
	checkpoint.PartialResultsOfficine[gov][del] = officines

	return nil
}

func extractOfficines(page *rod.Page) ([]*Officine, error) {
	officines := make([]*Officine, 0)

	// Check if the page contains "Aucune officine installée dans cette zone"
	body, err := page.Element("body")
	if err != nil {
		return nil, fmt.Errorf("failed to get body element: %w", err)
	}

	font, err := body.Element("font")
	if err == nil { // If font exists, check its text
		text, _ := font.Text()
		if strings.Contains(text, "ucune officin") {
			return officines, nil // Return empty list
		}
	}

	// Find the last table in the body
	tables, err := page.Elements("table")
	if err != nil || len(tables) == 0 {
		return nil, fmt.Errorf("failed to find tables: %w", err)
	}

	lastTable := tables[len(tables)-1] // Select last table
	rows, err := lastTable.Elements("tr")
	if err != nil {
		return nil, fmt.Errorf("failed to get table rows: %w", err)
	}

	// Loop through table rows, skipping the header row
	for i := 1; i < len(rows); i++ {
		cells, err := rows[i].Elements("td")
		if err != nil || len(cells) < 3 {
			continue
		}

		officine := &Officine{
			Ordre: fmt.Sprint(i),
		}

		// Assign values based on column index
		for j, cell := range cells {
			text, _ := cell.Text()

			switch j {
			case 0:
				officine.Nom = text
			case 1:
				officine.Adresse = text
			case 2:
				officine.Telephone = text
			}
		}

		officines = append(officines, officine)
	}

	return officines, nil
}
