package mc

import (
	"bufio"
	"bytes"
	"encoding/json"
)

type ServerStatus struct {
	Version struct {
		Name string			`json:"name"`
		Protocol int		`json:"protocol"`
	}						`json:"version"`

	Players struct {
		Max int				`json:"max"`
		Online int			`json:"online"`
		Sample []struct {
			Name string		`json:"name"`
			Id string		`json:"id"`
		}					`json:"sample"`
	}						`json:"players"`

	Description struct {
		Text string			`json:"text"`
	}						`json:"description"`

	Favicon string			`json:"favicon"`
}

/* Response packets are identical for any given server status, 
 * so we should save resources by not creating a new one
 * for every request. */
var responsePackets = map[*ServerStatus]*Packet{}

func DecodePacketPing(p *Packet) (int64, error) {
	r := Reader{bufio.NewReader(bytes.NewReader(p.Data))}

	ping, err := r.ReadInt64()

	if err != nil {
		return -1, err
	}
	return ping, nil
}

func NewPacketPing(ping int64) (*Packet, error){
	p := Packet{Id: 0x01}
	buf := new(bytes.Buffer)
	w := Writer{bufio.NewWriter(buf)}
	
	err := w.WriteInt64(ping)
	if err != nil {
		return nil, err
	}

	w.Flush()

	p.Data = buf.Bytes()
	return &p, nil
}

func NewPacketResponse(status *ServerStatus) (*Packet, error){
	if responsePackets[status] == nil {
		responsePackets[status] = &Packet{Id: 0x00}
		buf := new(bytes.Buffer)
		w := Writer{bufio.NewWriter(buf)}

		bs, err := json.Marshal(status)
		if err != nil {
			return nil, err
		}

		err = w.WriteVarstring(string(bs))
		if err != nil {
			return nil, err
		}

		w.Flush()

		responsePackets[status].Data = buf.Bytes()
	}

	return responsePackets[status], nil
}
