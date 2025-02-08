package src

import (
	"flag"

	"github.com/charmbracelet/log"
)

type Flags struct {
	Debug bool
	Mode  uint8
	Time  uint8
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
	debug := flag.Bool("debug", false, "Enable debug mode")
	mode := flag.String("mode", "attente", "Mode: attente OR officine")
	time := flag.String("time", "jour", "Time: jour OR nuit")
	flag.Parse()

	if *mode != "attente" && *mode != "officine" {
		log.Fatal("Invalid mode, must be 'attente' or 'officine'")
	}

	if *time != "jour" && *time != "nuit" {
		log.Fatal("Invalid time, must be 'jour' or 'nuit'")
	}

	flags := &Flags{
		Debug: *debug,
	}

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
