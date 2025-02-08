package src

const (
	LISTE_ATTENTE = "http://www.dpm.tn/dpm_pharm/asppharm/listgouv_att.php"
	OFFICINE      = "http://www.dpm.tn/dpm_pharm/asppharm/listgouv.php"
	JOUR          = "ON"
	NUIT          = "OFF"
	RESULTS_DIR   = "results"
)

func checkpointFileName(flags *Flags) string {
	var name string

	switch flags.Mode {
	case MODE_ATTENTE:
		name = "liste_attente"
	case MODE_OFFICINE:
		name = "liste_officine"
	}

	switch flags.Time {
	case TIME_JOUR:
		name += "_jour"
	case TIME_NUIT:
		name += "_nuit"
	}
	name += ".checkpoint.json"

	return name
}

func resultsFileName(flags *Flags) string {
	var name string

	switch flags.Mode {
	case MODE_ATTENTE:
		name = "liste_attente"
	case MODE_OFFICINE:
		name = "liste_officine"
	}

	switch flags.Time {
	case TIME_JOUR:
		name += "_jour"
	case TIME_NUIT:
		name += "_nuit"
	}

	name += ".json"

	return name
}
