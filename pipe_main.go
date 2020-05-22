package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

var ongoingDiagProcess bool

// Diagnostic code
func pipeMain() {
	log.Println("Ouverture de la pipe nommée: locuste.brain.diagnostic")
	readPipe, _ := os.OpenFile("locuste.brain.diagnostic", os.O_RDONLY, 0600)
	defer readPipe.Close()
	writePipe, _ := os.OpenFile("locuste.diagnostic.brain", os.O_WRONLY, 0600)
	defer writePipe.Close()

	var buffer bytes.Buffer
	log.Println("En attente d'instruction de diagnostiques")
	ongoingDiagProcess = true

	for ongoingDiagProcess == true {
		_, err := io.Copy(&buffer, readPipe)

		if buffer.Len() > 0 {
			log.Println("Instruction de diagnostique ", buffer.String())
		}

		if err != nil {
			if err.Error() != "EOF" {
				return
			}
		}

		time.Sleep(250 * time.Millisecond)
	}

	log.Println("Arrêt du module de diagnostique")

}
