package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/fysac/golem/mc"
)

var State = map[net.Conn]int{}
var StateLock = sync.Mutex{}

func handleConn(c net.Conn) {
	defer func() {
		c.Close()

		StateLock.Lock()
		delete(State, c)
		StateLock.Unlock()
	}()

	log.Println("Connection from", c.RemoteAddr())

	r := mc.Reader{bufio.NewReader(c)}
	w := mc.Writer{bufio.NewWriter(c)}

	for {
		p, err := r.ReadPacket()
		if err != nil {
			if err != io.EOF {
				log.Println("Error fetching packet:", err)
			}
			return
		}

		err = handlePacket(p, &w, c)
		if err != nil {
			if err != io.EOF {
				log.Println("Error decoding packet:", err)
			}
			return
		}
	}
}

func handlePacket(p *mc.Packet, w *mc.Writer, c net.Conn) error {
	switch State[c] {
	case mc.StateHandshake:
		if p.Id == mc.PacketHandshakeId {
			version, addr, port, state, err := mc.DecodePacketHandshake(p)
			if err != nil {
				return err
			}

			StateLock.Lock()
			State[c] = state
			StateLock.Unlock()

			log.Println("Handshake:", version, addr, port, state)
		}

	case mc.StateStatus:
		if p.Id == mc.PacketStatusRequestId {
			log.Println("Request")

			p, err := mc.NewPacketResponse(&Status)
			if err != nil {
				return err
			}

			return w.WritePacket(p)

		} else if p.Id == mc.PacketPingId {
			ping, err := mc.DecodePacketPing(p)
			if err != nil {
				return err
			}

			p, err := mc.NewPacketPing(ping)
			if err != nil {
				return err
			}

			return w.WritePacket(p)
		}

	case mc.StateLogin:
		if p.Id == mc.PacketLoginStartId {
			username, err := mc.DecodePacketLoginStart(p)
			if err != nil {
				return err
			}

			log.Println("Login start:", username)

			p, err := mc.NewPacketDisconnectLogin("Not implemented, " + username + "!")
			if err != nil {
				return err
			}

			return w.WritePacket(p)
		}

	case mc.StatePlay:
		log.Println("State play not implemented")

	default:
		return fmt.Errorf("Connection %v in invalid state %v", c.RemoteAddr(), State[c])
	}

	return nil
}
