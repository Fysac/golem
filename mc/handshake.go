package mc

import (
	"bufio"
	"bytes"
)

func DecodePacketHandshake(p *Packet) (int, string, uint16, int, error){
	r := Reader{bufio.NewReader(bytes.NewReader(p.Data))}

	version, err := r.ReadVarint()
	if err != nil {
		return -1, "", 0, -1, err
	}

	addr, err := r.ReadVarstring()
	if err != nil {
		return -1, "", 0, -1, err
	}

	port, err := r.ReadUint16()
	if err != nil {
		return -1, "", 0, -1, err
	}

	state, err := r.ReadVarint()
	if err != nil {
		return -1, "", 0, -1, err
	}

	return version, addr, port, state, nil
}
