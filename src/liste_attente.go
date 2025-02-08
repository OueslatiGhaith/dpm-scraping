package src

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
)

func scrapeWaitingList(ctx context.Context, page *rod.Page, checkpoint *Checkpoint, flags *Flags) error {
	log.Infof("Navigate to %s", LISTE_ATTENTE)
	if err := page.Navigate(LISTE_ATTENTE); err != nil {
		return fmt.Errorf("failed to navigate to %s: %w", LISTE_ATTENTE, err)
	}

	log.Debug("Getting list of gouvernourats")
	govSelect := page.MustElement("select[name='cod_gouv']")
	options := govSelect.MustElements("option")
	log.Debugf("Got %d gouvernourats", len(options))

	var gouvernourats []string
	for _, opt := range options {
		gov := opt.MustText()
		gouvernourats = append(gouvernourats, gov)
	}

	for _, gov := range gouvernourats {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
		}

		// skip if already processed
		if _, ok := checkpoint.ProcessedGovs[gov]; ok {
			log.Warnf("\tGouvernourat %s already processed, skipping", gov)
			continue
		}

		log.Infof("Processing gouvernourat %s", gov)
		checkpoint.CurrentGov = gov

		if err := processGouvernourat(ctx, page, gov, checkpoint, flags); err != nil {
			return fmt.Errorf("failed to process gouvernourat %s: %w", gov, err)
		}

		// mark gouvernourat as processed
		checkpoint.ProcessedGovs[gov] = true
		checkpoint.CurrentGov = ""
		checkpoint.CurrentDel = ""

		// save checkpoint after each government
		if err := saveCheckpoint(checkpoint, flags); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	return nil
}

func processGouvernourat(ctx context.Context, page *rod.Page, gov string, checkpoint *Checkpoint, flags *Flags) error {
	// navigate to the main page and select the gouvernourat
	if err := navigateToGouvernourat(page, gov, flags); err != nil {
		return fmt.Errorf("failed to navigate to gouvernourat: %w", err)
	}

	// waiting for delegations select to be ready
	delSelect := page.MustElement("select[name='cod_del']")
	options := delSelect.MustElements("option")

	log.Debug("Getting list of delegations")
	var delegations []string
	for _, opt := range options {
		del := opt.MustText()
		delegations = append(delegations, del)
	}

	// initialize map for this gouvernourat if needed
	if checkpoint.PartialResultsAttente[gov] == nil {
		checkpoint.PartialResultsAttente[gov] = make(map[string]*Attente)
	}

	// process each delegation
	for _, del := range delegations {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
		}

		// skip if already processed
		if slices.Contains(checkpoint.ProcessedDels[gov], del) {
			log.Warnf("Delegation %s already processed, skipping", del)
			continue
		}

		log.Infof("\tProcessing delegation %s", del)
		if err := processDelegation(page, gov, del, checkpoint, flags); err != nil {
			if saveErr := saveCheckpoint(checkpoint, flags); saveErr != nil {
				log.Errorf("Failed to save checkpoint: %s", saveErr)
			}
			return fmt.Errorf("failed to process delegation %s: %w", del, err)
		}

		// update checkpoint
		if checkpoint.ProcessedDels[gov] == nil {
			checkpoint.ProcessedDels[gov] = make([]string, 0)
		}
		checkpoint.ProcessedDels[gov] = append(checkpoint.ProcessedDels[gov], del)

		// save checkpoint after each delegation
		if err := saveCheckpoint(checkpoint, flags); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	return nil
}

func processDelegation(page *rod.Page, gov, del string, checkpoint *Checkpoint, flags *Flags) error {
	// navigate to the main page and select the gouvernourat
	if err := navigateToGouvernourat(page, gov, flags); err != nil {
		return fmt.Errorf("failed to navigate to gouvernourat: %w", err)
	}

	// select delegation
	log.Debugf("Selecting delegation %s", del)
	page.MustElement("select[name='cod_del']").MustSelect(del)
	page.MustElement("input[type='submit']").MustClick()

	// wait for the page to be ready
	log.Debugf("Waiting for page to be ready")
	page.MustWaitLoad()

	attente, err := extractAttente(page)
	if err != nil {
		return fmt.Errorf("failed to extract attente: %w", err)
	}

	// store in checkpoint
	if checkpoint.PartialResultsAttente[gov] == nil {
		checkpoint.PartialResultsAttente[gov] = make(map[string]*Attente)
	}
	checkpoint.PartialResultsAttente[gov][del] = attente

	return nil
}

func extractAttente(page *rod.Page) (*Attente, error) {
	log.Debug("Extracting attente")

	attente := &Attente{}

	// check for empty state
	fonts, err := page.Elements("font[color='#000000']")
	if err != nil {
		return nil, fmt.Errorf("failed to get font elements: %w", err)
	}

	log.Debugf("Got %d font elements", len(fonts))
	for _, font := range fonts {
		text, err := font.Text()
		if err != nil {
			continue
		}
		if strings.TrimSpace(text) == "Pas d'inscription sur cette liste" {
			return attente, nil
		}
	}

	// extract zone, population, nb_officines using XPath
	attente.Zone = getByXPath(page, `/html/body/table/tbody/tr[1]/td/p/font/b/b/font[2]`)
	log.Debug("Extracted zone")

	attente.Population = getByXPath(page, "/html/body/table/tbody/tr[1]/td/p/font/b/b/b/font[2]")
	log.Debug("Extracted population")

	attente.NbOfficines = getByXPath(page, "/html/body/table/tbody/tr[1]/td/p/font/b/b/b/b/font[2]")
	log.Debug("Extracted nbOfficines")

	// get waiting list from last table
	table, err := page.Elements("table:last-child")
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	if len(table) == 0 {
		return nil, fmt.Errorf("no table found")
	}
	log.Debugf("Got %d tables", len(table))

	lastTable := table[len(table)-1]
	rows, err := lastTable.Elements("tr")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	// skip header row
	for i := 1; i < len(rows); i++ {
		cells, err := rows[i].Elements("td")
		if err != nil {
			continue
		}

		person := PersonneAttente{}
		for j, cell := range cells {
			child, err := cell.Element("*")
			if err != nil {
				continue
			}

			text, err := child.Text()
			if err != nil {
				continue
			}

			switch j {
			case 0:
				person.Ordre = text
			case 1:
				person.Nom = text
			case 2:
				person.Prenom = text
			case 3:
				person.Epouse = text
			case 4:
				person.DateInscription = text
			}
		}

		attente.Liste = append(attente.Liste, person)
	}

	return attente, nil
}

func navigateToGouvernourat(page *rod.Page, gouvernourat string, flags *Flags) error {
	// navigate to the main page and select the gouvernourat
	log.Debugf("Navigating to gouvernourat %s", gouvernourat)
	if err := page.Navigate(LISTE_ATTENTE); err != nil {
		return fmt.Errorf("failed to navigate to main page: %w", err)
	}

	// select gouvernourat and JOUR option
	log.Debugf("Selecting gouvernourat %s", gouvernourat)
	page.MustElement("select[name='cod_gouv']").MustSelect(gouvernourat)

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
