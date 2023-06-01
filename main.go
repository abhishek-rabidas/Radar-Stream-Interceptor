package main

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

var START = []byte{0xca, 0xcb, 0xcc, 0xcd}
var END = []byte{0xea, 0xeb, 0xec, 0xed}

type Radar struct {
	Logger *Logging
}

func main() {
	config, err := LoadConfig()

	var radar Radar = Radar{Logger: Init(config.OutputDir)}

	if err != nil {
		panic(err)
	}

	address := config.IP + ":" + config.Port

	raddr, err := net.ResolveTCPAddr("tcp", address)

	connection, err := net.DialTCP("tcp", nil, raddr)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("-----Connection established [%s]-----", address)

	radar.ReadData(connection)
}

func (r *Radar) ReadData(stream io.ReadWriteCloser) {
	for {
		packet, err := readPacket(stream)
		if err == nil {
			r.HandlePacket(packet)
		}
	}
}

func readPacket(stream io.ReadWriteCloser) ([]byte, error) {
	buf := make([]byte, 1)
	//start market
	i := 0
	for {
		_, err := stream.Read(buf)
		if err == io.EOF {
			return nil, err
		}
		//log.Print(cnt, buf[0], START[i])
		if buf[0] == START[i] {
			//log.Print("M")
			i++
		} else {
			//start matching from start
			i = 0
		}
		//header is matched
		if i == len(START) {
			break
		}
	}
	packet := make([]byte, 512)
	i = 0
	cnt := 0
	for {
		_, err := stream.Read(buf)
		if err == io.EOF {
			return nil, err
		}
		//log.Print(cnt, buf[0], START[i])
		if buf[0] == END[i] {
			//log.Print("M")
			i++
		} else {
			packet[cnt] = buf[0]
			cnt++
			//start matching from start
			i = 0
		}
		//header is matched
		if i == len(END) {
			return packet[0:cnt], nil
		}
	}
}

func reverse(payload []byte) []byte {
	newpayload := make([]byte, len(payload))
	for idx := range payload {
		newpayload[idx] = payload[len(payload)-idx-1]
	}
	return newpayload
}

func (r *Radar) HandlePacket(packet []byte) {
	idx := 0
	packetSize := len(packet)
	for {
		msgType := binary.BigEndian.Uint16(packet[idx : idx+2])
		idx += 2
		msgLen := int(packet[idx])
		idx += 1
		payload := packet[idx : idx+msgLen]
		idx += msgLen
		if idx >= packetSize-1 {
			break
		}
		payload = reverse(payload)
		switch msgType {
		case 0x0500:
			HandleStatusMessage(payload)
		case 0x0501:
			HandleObjectMessage(payload)
		case 0x02ff:
			HandleSyncMessage(payload)
		case 0x0734:
			HandleStatusMessage(payload)
		default:
			if msgType >= 0x0502 && msgType <= 0x057F {
				log.Info("Detection MESSAGE")
				r.Logger.WriteObjectDetectionData(HandleDetectionMessage(payload))
			} else {
				log.Infof("Unknown type: 0x%x\n", msgType)
			}
		}
	}
}
