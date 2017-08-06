package mc

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

const (
	MaxVarintSize   = 5
	MaxPacketSize   = 32767*4 + 3
	MaxStringLength = 32767*4 + 3
)

type Reader struct {
	*bufio.Reader
}

func (r *Reader) ReadPacket() (*Packet, error) {
	len, err := r.ReadVarint()
	if err != nil {
		return nil, err
	}

	if len <= 0 {
		return nil, fmt.Errorf("Invalid packet size: 0x%02x", len)
	}
	if len > MaxPacketSize {
		return nil, fmt.Errorf("Packet too big")
	}

	id, err := r.ReadVarint()
	if err != nil {
		return nil, err
	}

	if id < 0 || id > 0xff {
		return nil, fmt.Errorf("Invalid id: 0x%02x", id)
	}

	data := make([]byte, len-1)

	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}

	p := Packet{Length: len, Id: id, Data: data}
	log.Printf("Received packet: id=0x%02x, len=0x%02x, data=%#02v", p.Id, p.Length, p.Data)
	return &p, nil
}

/* Read a long */
func (r *Reader) ReadInt64() (int64, error) {
	const size = 8
	var res int64 = 0

	bs := make([]byte, size)

	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}

	for i := range bs {
		res |= int64(bs[i]) << uint(8*(len(bs)-i-1))
	}

	return res, nil
}

/* Read an unsigned short */
func (r *Reader) ReadUint16() (uint16, error) {
	const size = 2
	var res uint16 = 0

	bs := make([]byte, size)

	_, err := r.Read(bs)
	if err != nil {
		return 0, err
	}

	for i := range bs {
		res |= uint16(bs[i]) << uint(8*(len(bs)-i-1))
	}

	return res, nil
}

/* Read a Minecraft varint */
func (r *Reader) ReadVarint() (int, error) {
	var i uint = 0
	var b byte = 0x80
	var err error
	res := 0

	for b&0x80 != 0 {
		b, err = r.ReadByte()
		if err != nil {
			return -1, err
		}

		res |= (int(b&byte(0x7f)) << (7 * i))
		i++

		if i > MaxVarintSize {
			return -1, fmt.Errorf("Varint too big")
		}
	}
	return res, nil
}

/* Read a Minecraft string */
func (r *Reader) ReadVarstring() (string, error) {
	len, err := r.ReadVarint()
	if err != nil {
		return "", err
	}

	if len <= 0 || len > MaxStringLength {
		return "", fmt.Errorf("Invalid string length: %v", len)
	}

	bs := make([]byte, len)

	_, err = r.Read(bs)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

type Writer struct {
	*bufio.Writer
}

func (w *Writer) WritePacket(p *Packet) error {
	p.Length = len(p.Data) + 1

	buf := new(bytes.Buffer)
	bufw := Writer{bufio.NewWriter(buf)}

	err := bufw.WriteVarint(p.Id)
	if err != nil {
		return err
	}

	_, err = bufw.Write(p.Data)
	if err != nil {
		return err
	}

	bufw.Flush()

	err = w.WriteVarint(p.Length)
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	w.Flush()
	log.Printf("Sent packet: id=0x%02x, len=0x%02x, data=%#02v", p.Id, p.Length, p.Data)
	return nil
}

/* Write a long */
func (w *Writer) WriteInt64(value int64) error {
	const size = 8
	bs := make([]byte, size)

	for i := range bs {
		bs[i] = byte(value & 0xff)
		value >>= 8

		if value == 0 {
			break
		}
	}

	for i := len(bs) - 1; i >= 0; i-- {
		err := w.WriteByte(bs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

/* Write a Minecraft varint */
func (w *Writer) WriteVarint(value int) error {
	for {
		temp := byte(value & 0x7f)
		value >>= 7

		if value != 0 {
			temp |= 0x80
		}

		err := w.WriteByte(temp)
		if err != nil {
			return err
		}

		if value == 0 {
			break
		}
	}
	return nil
}

/* Write a Minecraft string */
func (w *Writer) WriteVarstring(str string) error {
	err := w.WriteVarint(len([]byte(str)))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}
