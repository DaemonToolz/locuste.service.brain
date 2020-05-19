package main

import (
	"log"
	"os/exec"
	"fmt"
	"strings"
)


var rtspProcesses map[string]*exec.Cmd
var ffmpegProcesses map[string]*exec.Cmd
const (
	nodeCommand string ="node"
	nodeRtspPath string = "/home/pi/project/locuste/services/brain/node/rtsp-socket.js"

	ffmpegCommand string = "ffmpeg -f rtsp -rtsp_transport udp -i rtsp://%s/live -f mpegts -codec:v mpeg1video  -q 2 -an -b:v 50M -threads 32 -strict experimental -bf 0 -r 25 -muxdelay 0.001 http://localhost:%d/anafi"
)  

func initRtspListener(){
	rtspProcesses = make(map[string]*exec.Cmd)
	ffmpegProcesses = make(map[string]*exec.Cmd)
}


func startVideoServer(name string, port int){

	renew := false
	if process, ok := rtspProcesses[name]; ok {
		if process == nil || process.ProcessState != nil && process.ProcessState.Exited() || process.Process == nil {
			log.Println("Le processus JSMPEG-TS est HS pour ", name, " au port ", port, ". On redémarre")
			renew = true;
		} else {
			log.Println("Processus JSMPEG-TS déjà OK pour ", name, " au port ", port)
		}
	} else { renew = true }

	if renew {
		log.Println("Démarrage du processus de visulation pour ", name, " au port ", port)
		rtspProcesses[name] = exec.Command(string(nodeCommand), string(nodeRtspPath), "anafi", string(fmt.Sprintf("%d", port+50)), string(fmt.Sprintf("%d", port)))
		if err := rtspProcesses[name].Start(); err != nil {
			failOnError(err, "Le processus JSMPEG-TS n'a pas pu démarrer")
		}

		go func(processName string) { 
			defer func() {
				if r := recover(); r != nil {
					AddOrUpdateExtCompStatus(processName, VideoServer, false)
				}
			}()
			AddOrUpdateExtCompStatus(processName, VideoServer, true)
			err := rtspProcesses[processName].Wait()
			log.Println("Processus JSMPEG-TS arrêté pour ", name, " raison ", err)
			AddOrUpdateExtCompStatus(processName, VideoServer, false)
		}(name)
		log.Println("Processus JSMPEG-TS OK pour ", name, " at ", port)
	}

}

func startFfmpegStream(name string,port int){
	renew := false
	if process, ok := ffmpegProcesses[name]; ok {
		if process == nil || process.ProcessState != nil && process.ProcessState.Exited() || process.Process == nil {
			log.Println("Le processus RTSP->FFMPEG est HS pour ", name, " au port ", port, ". On redémarre")
			renew = true;
		} else {
			log.Println("Processus RTSP->FFMPEG déjà OK pour ", name, " au port ", port)
		}
	} else { renew = true }

	if renew {
		log.Println("Démarrage du processus de visulation pour ", name, " au port ", port)
		dip := name
		if index := strings.LastIndex(name, "_"); index != -1 {
			dip = name[index+1:]
		}
		command := strings.Fields(fmt.Sprintf(ffmpegCommand, dip, port+50))
		ffmpegProcesses[name] = exec.Command(command[0], command[1:]...)
		if err := ffmpegProcesses[name].Start(); err != nil {
			AddOrUpdateExtCompStatus(name, VideoStream, false)
			failOnError(err, "Le processus RTSP->FFMPEG n'a pas pu démarrer")
		}

		go func(processName string) { 
			defer func() {
				if r := recover(); r != nil {
					AddOrUpdateExtCompStatus(name, VideoStream, false)
				}
			}()
			AddOrUpdateExtCompStatus(name, VideoStream, true)
			err := ffmpegProcesses[processName].Wait()
			log.Println("Processus RTSP->FFMPEG arrêté pour ", name, " raison ", err)
			AddOrUpdateExtCompStatus(name, VideoStream, false)
		}(name)
		log.Println("Processus RTSP->FFMPEG OK pour ", name, " at ", port)
	}
}

func onClose(){
	for _, process := range rtspProcesses{
		process.Process.Kill()
	}
	for _, process := range ffmpegProcesses{
		process.Process.Kill()
	}
	
}