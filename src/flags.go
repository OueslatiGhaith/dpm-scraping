package src

import (
	"flag"
	"log"
)

type Flags struct {
	Mode uint8
	Time uint8
}

const (
	MODE_ATTENTE = iota
	MODE_OFFICINE
)

const (
	TIME_JOUR = iota
	TIME_NUIT
)

func parseFlags() *Flags {
	mode := flag.String("mode", "attente", "Mode: attente OR officine")
	time := flag.String("time", "jour", "Time: jour OR nuit")
	flag.Parse()

	if *mode != "attente" && *mode != "officine" {
		log.Fatal("Invalid mode, must be 'attente' or 'officine'")
	}

	if *time != "JOUR" && *time != "NUIT" {
		log.Fatal("Invalid time, must be 'jour' or 'nuit'")
	}

	flags := &Flags{}

	switch *mode {
	case "attente":
		flags.Mode = MODE_ATTENTE
	case "officine":
		flags.Mode = MODE_OFFICINE
	}

	switch *time {
	case "jour":
		flags.Time = TIME_JOUR
	case "nuit":
		flags.Time = TIME_NUIT
	}

	return flags
}
