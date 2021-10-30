package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

// Take a look at:
// encoding/binary.Write
// encoding/binary.BigEndian
// net.ListenIP

type ICMPMsg struct {
	Type     byte
	Code     byte
	Checkusm uint16
	Id       uint16
	SeqNum   uint16
	Payload  []byte
}

func main() {
	pkt := ICMPMsg{
		Type:     8,
		Code:     0,
		Checkusm: 0,
		Id:       5,
		SeqNum:   10,
		Payload:  []byte("Am I in an ICMP message O_o?"),
	}

	bpkt := pkt.Marshal()

	chksum := CheckSum(bpkt)

	pkt.Checkusm = chksum

	bpkt[2] = byte(chksum & (0xFF << 8) >> 8)
	bpkt[3] = byte(chksum)

	cnx, err := net.Dial("ip4:icmp", "127.0.0.1")
	if err != nil {
		fmt.Printf("Could not open a connection to 127.0.0.1: %v\n", err)
		os.Exit(1)
	}
	defer cnx.Close()

	sig_ch := make(chan os.Signal, 1)
	signal.Notify(sig_ch, os.Interrupt)

	go func() {
		for {
			cnx.Write(bpkt)
			time.Sleep(time.Second)
			fmt.Printf("Pinged localhost!\n")
		}
	}()

	<-sig_ch
	fmt.Printf("Goodbye :P\n")
	os.Exit(0)
}

func CheckSum(pkt []byte) uint16 {
	var tmp uint32
	// As we are append()ing data to a slice with no more capacity
	// we are just getting a pointer to a new slice. Thus, this
	// modification will NOT have an effect on the resulting packet.
	if len(pkt)%2 != 0 {
		pkt = append(pkt, 0)
	}
	for i := 0; i <= len(pkt)-2; i += 2 {
		x := uint32(pkt[i])<<8 | uint32(pkt[i+1])
		foo := tmp + x

		if foo > 0x10000 {
			foo++
		}
		tmp = foo & 0xFFFF
	}
	return uint16(^tmp & 0xFFFF)
}

func (p ICMPMsg) Marshal() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, p.Type)
	binary.Write(buf, binary.BigEndian, p.Code)
	binary.Write(buf, binary.BigEndian, p.Checkusm)
	binary.Write(buf, binary.BigEndian, p.Id)
	binary.Write(buf, binary.BigEndian, p.SeqNum)
	binary.Write(buf, binary.BigEndian, p.Payload)

	return buf.Bytes()
}
