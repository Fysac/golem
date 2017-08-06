package mc

import (
	"bufio"
	"bytes"
)

func DecodePacketLoginStart(p *Packet) (string, error) {
	r := Reader{bufio.NewReader(bytes.NewReader(p.Data))}

	username, err := r.ReadVarstring()
	if err != nil {
		return "", err
	}

	return username, nil
}

func NewPacketDisconnectLogin(reason string) (*Packet, error) {
	p := Packet{Id: 0x00}
	buf := new(bytes.Buffer)
	w := Writer{bufio.NewWriter(buf)}

	err := w.WriteVarstring("[\"" + reason + "\"]")
	if err != nil {
		return nil, err
	}

	w.Flush()

	p.Data = buf.Bytes()
	return &p, nil
}

func NewPacketLoginSuccess(username string, uuid string) (*Packet, error) {
	p := Packet{Id: 0x02}
	buf := new(bytes.Buffer)
	w := Writer{bufio.NewWriter(buf)}

	err := w.WriteVarstring(uuid)
	if err != nil {
		return nil, err
	}

	err = w.WriteVarstring(username)
	if err != nil {
		return nil, err
	}

	w.Flush()

	p.Data = buf.Bytes()
	return &p, nil
}
