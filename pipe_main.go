package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"
)

var ongoingDiagProcess bool

// Diagnostic code
func pipeMain() {
	log.Println("Ouverture de la pipe nommée pour locuste.service.brain")
	syscall.Mkfifo("/tmp/locuste.brain.diagnostic", 0666)
	syscall.Mkfifo("/tmp/locuste.diagnostic.brain", 0666)

	//defer readPipe.Close()
	//writePipe, _ := os.OpenFile("/tmp/locuste.diagnostic.brain", os.O_WRONLY, os.ModeNamedPipe)
	//defer writePipe.Close()

	log.Println("En attente d'instruction de diagnostiques")
	ongoingDiagProcess = true

	for ongoingDiagProcess == true {
		var buffer bytes.Buffer
		readPipe, _ := os.OpenFile("/tmp/locuste.brain.diagnostic", os.O_RDONLY, os.ModeNamedPipe)
		writePipe, _ := os.OpenFile("/tmp/locuste.diagnostic.brain", os.O_WRONLY, os.ModeNamedPipe)
		_, err := io.Copy(&buffer, readPipe)
		if buffer.Len() > 0 {
			log.Println("Instruction de diagnostique ", buffer.String())
			writePipe.WriteString(fmt.Sprintf("Instruction %s recue", buffer.String()))
			log.Println("Réponse envoyée")
		}

		if err != nil {
			failOnError(err, "Une erreur est survenue")
			if err.Error() != "EOF" {
				readPipe.Close()
				break
			}
		}
		readPipe.Close()
		writePipe.Close()
		time.Sleep(250 * time.Millisecond)
	}

	log.Println("Arrêt du module de diagnostique")

}
