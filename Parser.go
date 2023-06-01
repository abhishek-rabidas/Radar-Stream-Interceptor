package main

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

func HandleDetectionMessage(payload []byte) *ObjectDetection {
	var object ObjectDetection = ObjectDetection{}
	value := binary.LittleEndian.Uint64(payload[0:8])
	xp1 := value & 0b00000000_00000000_00000000_00000000_00000000_00000000_00111111_11111110 >> 1  //1..13
	yp1 := value & 0b00000000_00000000_00000000_00000000_00000111_11111111_11000000_00000000 >> 14 //14..26
	sx1 := value & 0b00000000_00000000_00000000_00111111_11111000_00000000_00000000_00000000 >> 27 //27..37
	sy1 := value & 0b00000000_00000001_11111111_11000000_00000000_00000000_00000000_00000000 >> 38 //38..48
	oll := value & 0b00000000_11111110_00000000_00000000_00000000_00000000_00000000_00000000 >> 49 //49..55
	oid := value & 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000 >> 56 //56..63

	object.ObjectID = int(oid)
	object.X = (float32(xp1) - 4096) * 0.128
	object.Y = (float32(yp1) - 4096) * 0.128
	object.X_Speed = (float32(sx1) - 1024) * 0.1
	object.Y_Speed = (float32(sy1) - 1024) * 0.1
	object.Length = float32(oll) * 0.2
	return &object
	//log.Infof("DET MSG>SIG: %d, %d, xp1: %.2f, yp1 %.2f, sx: %.2f, sy: %.2f,ol : %.2f \n", oid, signal, xp, yp, sx, sy, objlen)
}

func HandleStatusMessage(payload []byte) {
	//status := payload[0]
	//imni := payload[1]
	//split after masking
	//im := 0x0f & imni
	//ni := (0xf0 & imni) >> 4
	//diagnose := payload[2]
	//reserve := payload[3]
	//ts := binary.LittleEndian.Uint32(payload[4:8])
	//log.Infof("STATUS MSG>TS: %d, status: %d", ts, status)
}

func HandleObjectMessage(payload []byte) {
	numObjects := payload[0]
	numMsgs := payload[1]
	cycleDuration := payload[2]
	//odf := payload[3]
	//odf0 := 0x0f & odf
	//odf1 := (0xf0 & odf) >> 4
	cyclecount := binary.LittleEndian.Uint32(payload[4:8])
	log.Infof("OBJ MSG>CC: %d, dur: %d, numObjects %d, numMsgs %d\n", cyclecount, cycleDuration, numObjects, numMsgs)
}

func HandleSyncMessage(payload []byte) {
	//ts := binary.LittleEndian.Uint32(payload[0:4])
	//log.Info("SYNC MSG>TS:", ts)
}
