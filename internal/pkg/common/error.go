package common

import "log"

func LogError(errors ...error) {
	for _, err := range errors {
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}
	}
}
