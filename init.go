package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//9010
//192.168.1.67

var serveMux *http.ServeMux

func main() {
	initModuleRestartMapper()
	initOperatorDictionary()
	initDroneConfiguration()
	initDroneSettings()
	initHealthMonitor()

	AddOrUpdateStatus(BrainMainRunner, true)
	initConfiguration()
	prepareLogs()
	go pipeMain()
	initRtspListener()
	RestartRPCServer()
	initFileWatcher("/home/pi/project/locuste/data") // retirer les chemins d'accès
	initSocketServer()

	initRPCClient()

	serveMux = http.NewServeMux()
	router = NewRouter()
	initMiddleware(router)

	RestartHTTPServer()

	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		log.Println("WebSocket Server online")
		serveMux.Handle("/socket.io/", server)
	}()

	RestartSocketServer()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)

	select {
	case <-sigChan:
		ongoingDiagProcess = false
		pulse.Stop()
		go func() { stopCondition <- true }()
		AddOrUpdateStatus(BrainSocketServer, false)
		AddOrUpdateStatus(BrainHttpServer, false)
		AddOrUpdateStatus(BrainMainRunner, false)
		AddOrUpdateStatus(BrainWatcher, false)
		AddOrUpdateStatus(BrainRPCServer, false)
		Unregister()

		time.Sleep(5 * time.Second)
		onClose()
		watcher.Close()
		logFile.Close()
		os.Exit(0)
	}
}

// RestartSocketServer Redémarrage du serveur / module de WebSocket
func RestartSocketServer() {
	if result, ok := GlobalStatuses[BrainSocketServer]; !ok || (ok && !result) {
		initConfiguration()
		go func() {
			defer func() {
				if r := recover(); r != nil {
					AddOrUpdateStatus(BrainSocketServer, false)
				}
			}()
			AddOrUpdateStatus(BrainSocketServer, true)
			log.Println("Serving Websocket at ", appConfig.socketListenUri(), "")
			log.Println(http.ListenAndServe(appConfig.socketListenUri(), serveMux))
			AddOrUpdateStatus(BrainSocketServer, false)
		}()
	}
}

// RestartHTTPServer Redémarrage du serveur / module HTTP
func RestartHTTPServer() {
	if result, ok := GlobalStatuses[BrainHttpServer]; !ok || (ok && !result) {
		initConfiguration()
		go func() {
			defer func() {
				if r := recover(); r != nil {
					AddOrUpdateStatus(BrainHttpServer, false)
				}
			}()

			AddOrUpdateStatus(BrainHttpServer, true)
			log.Println("Serving at ", appConfig.httpListenUri(), "")
			log.Println(http.ListenAndServe(appConfig.httpListenUri(), router))
			AddOrUpdateStatus(BrainHttpServer, false)
		}()
	}
}
