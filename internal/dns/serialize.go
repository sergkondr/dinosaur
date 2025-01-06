package dns

import (
	"encoding/binary"
)

const (
	RECORD_TYPE_A   = 1
	RECORD_CLASS_IN = 1
)

func (m *Message) Serialize() ([]byte, error) {
	// 12 bytes is a minimal length of a DNS message
	// because the message header itself has a length == 12 bytes
	data := make([]byte, 12)

	binary.BigEndian.PutUint16(data[0:2], m.Header.ID)
	binary.BigEndian.PutUint16(data[2:4], serializeFlags(m.Header.Flags))
	binary.BigEndian.PutUint16(data[4:6], uint16(1)) //
	binary.BigEndian.PutUint16(data[6:8], uint16(len(m.Answer)))

	queryData := serializeQuery(m.Question)
	data = append(data, queryData...)

	answerData := serializeAnswer(m.Answer)
	data = append(data, answerData...)

	return data, nil
}

func serializeFlags(flags Flags) uint16 {
	var data uint16

	data = uint16(flags.Query) << 15
	data = data | uint16(flags.Opcode&0xF)<<11
	data = data | uint16(flags.AuthoritativeAnswer)<<10
	data = data | uint16(flags.Truncation)<<9
	data = data | uint16(flags.RecursionDesired)<<8
	data = data | uint16(flags.RecursionAvailable)<<7
	data = data | uint16(flags.Z&0x7)<<4
	data = data | uint16(flags.ResponseCode&0xF)

	return data
}

func serializeQuery(q MessageQuestion) []byte {
	data := make([]byte, 0)

	for i := 0; i < len(q.Qname); i++ {
		data = append(data, uint8(len(q.Qname[i])))
		for j := 0; j < len(q.Qname[i]); j++ {
			data = append(data, q.Qname[i][j])
		}
	}
	data = append(data, uint8(0)) // Splitter
	data = binary.BigEndian.AppendUint16(data, q.Qtype)
	data = binary.BigEndian.AppendUint16(data, q.Qclass)

	return data
}

func serializeAnswer(a []MessageResourceRecord) []byte {
	data := make([]byte, 0)

	// todo: message compression:
	// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
	for recordNumber := 0; recordNumber < len(a); recordNumber++ {
		data = append(data, []byte{0xc0, 0xc}...) // Compression pointer

		data = binary.BigEndian.AppendUint16(data, a[recordNumber].Type)
		data = binary.BigEndian.AppendUint16(data, a[recordNumber].Class)
		data = binary.BigEndian.AppendUint32(data, a[recordNumber].TTL)

		data = binary.BigEndian.AppendUint16(data, uint16(len(a[recordNumber].Data)))
		for octet := 0; octet < len(a[recordNumber].Data); octet++ {
			data = append(data, uint8(a[recordNumber].Data[octet]))
		}
	}

	return data
}
