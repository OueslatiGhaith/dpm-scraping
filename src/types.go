package src

import "time"

type Attente struct {
	Zone        string            `json:"zone"`
	Population  string            `json:"population"`
	NbOfficines string            `json:"nb_officines"`
	Liste       []PersonneAttente `json:"liste"`
}

type PersonneAttente struct {
	Ordre           string `json:"ordre"`
	Nom             string `json:"nom"`
	Prenom          string `json:"prenom"`
	Epouse          string `json:"epouse"`
	DateInscription string `json:"date_inscription"`
}

type ListeAttente map[string]map[string]*Attente

type Checkpoint struct {
	LastUpdated    time.Time           `json:"last_updated"`
	CurrentGov     string              `json:"current_gov"`
	CurrentDel     string              `json:"current_del"`
	ProcessedGovs  map[string]bool     `json:"processed_govs"`
	ProcessedDels  map[string][]string `json:"processed_dels"`
	PartialResults ListeAttente        `json:"partial_results"`
}
