package mc

type Packet struct {
	Length int
	Id int
	Data []byte
}

/* Game states:
 * http://wiki.vg/How_to_Write_a_Server */
const (
	StateHandshake = 0
	StateStatus = 1
	StateLogin = 2
	StatePlay = 3
)

/* Packet IDs:
 * http://wiki.vg/Protocol */
const (
	PacketHandshakeId = 0x00
	PacketStatusRequestId = 0x00
	PacketLoginStartId = 0x00
	PacketPingId = 0x01
)
