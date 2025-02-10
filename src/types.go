package src

import "time"

type (
	Gouvernourat string
	Delegation   string
)

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

type Officine struct {
	Ordre     string `json:"ordre"`
	Nom       string `json:"nom"`
	Adresse   string `json:"adresse"`
	Telephone string `json:"telephone"`
}

type Checkpoint struct {
	LastUpdated            time.Time                                   `json:"last_updated"`
	CurrentGov             Gouvernourat                                `json:"current_gov"`
	CurrentDel             Delegation                                  `json:"current_del"`
	ProcessedGovs          map[Gouvernourat]bool                       `json:"processed_govs"`
	ProcessedDels          map[Gouvernourat][]Delegation               `json:"processed_dels"`
	PartialResultsAttente  map[Gouvernourat]map[Delegation]*Attente    `json:"partial_results"`
	PartialResultsOfficine map[Gouvernourat]map[Delegation][]*Officine `json:"partial_results_officine"`
}
