package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Logging struct {
	OutputDir  string
	OutputPath string
	OutputFile *os.File
}

func Init(dir string) *Logging {
	folderPath := dir
	_, err := os.Stat(folderPath)

	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Output Folder Doesn't Exist, creating new one!")
			err = os.Mkdir("output", 0755)
			if err != nil {
				log.Error(err)
			}
		} else {
			log.Error(err)
		}
	}

	outputPath := "./output/log_" + time.Now().Format("2006-01-02 15:04:05") + ".log"

	outputPath = strings.ReplaceAll(outputPath, ":", "")

	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err)
	}

	return &Logging{
		OutputDir:  folderPath,
		OutputPath: outputPath,
		OutputFile: file,
	}

}

func (l *Logging) WriteObjectDetectionData(data *ObjectDetection) {
	bytesData, _ := json.Marshal(data)
	l.OutputFile.Write(bytesData)
	l.OutputFile.WriteString("\n")
}
