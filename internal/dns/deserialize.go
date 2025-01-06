package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

func Deserialize(data []byte) (*Message, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("data is too short")
	}

	msg := Message{}

	msg.Header.ID = binary.BigEndian.Uint16(data[0:2])

	flags, err := parseFlags(data[2:4])
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}
	msg.Header.Flags = flags

	msg.Header.QueriesCount = binary.BigEndian.Uint16(data[4:6])
	msg.Header.AnswerCount = binary.BigEndian.Uint16(data[6:8])
	msg.Header.NameServerCount = binary.BigEndian.Uint16(data[8:10])
	msg.Header.AdditionalCount = binary.BigEndian.Uint16(data[10:12])

	if err := msg.parseDNSNameInQuery(data[12:]); err != nil {
		return nil, fmt.Errorf("error parsing DNS names: %w", err)
	}

	return &msg, nil
}

func parseFlags(data []byte) (Flags, error) {
	if len(data) != 2 {
		return Flags{}, errors.New("wrong data length")
	}

	fl := Flags{}
	flags := binary.BigEndian.Uint16(data[0:2])
	fl.Query = uint8(flags >> 15)
	fl.Opcode = uint8((flags >> 11) & 0xF)
	fl.AuthoritativeAnswer = uint8((flags >> 10) & 1)
	fl.Truncation = uint8((flags >> 9) & 1)
	fl.RecursionDesired = uint8((flags >> 8) & 1)
	fl.RecursionAvailable = uint8((flags >> 7) & 1)
	fl.Z = uint8((flags >> 4) & 0x7)
	fl.ResponseCode = uint8(flags & 0xF)

	return fl, nil
}

func (m *Message) parseDNSNameInQuery(data []byte) error {
	dnsName := make([]string, 0)
	for {
		l := int(data[0])
		if l == 0 {
			break
		}
		dnsName = append(dnsName, string(data[1:l+1]))
		data = data[l+1:]
	}
	data = data[1:] // skip separator
	m.Question.Qname = dnsName
	m.Question.Qtype = binary.BigEndian.Uint16(data[0:2])
	m.Question.Qclass = binary.BigEndian.Uint16(data[2:4])

	return nil
}

func (m *Message) String() string {
	var str strings.Builder

	str.WriteString("DNS ")
	if m.Header.Flags.Query == 1 {
		str.WriteString("response")
	} else {
		switch m.Header.Flags.Opcode {
		case 0:
			str.WriteString("query")
		case 1:
			str.WriteString("inverse query")
		case 2:
			str.WriteString("server status")
		case 3:
			str.WriteString("notify")
		case 4:
			str.WriteString("update")
		default:
			str.WriteString("unknown")
		}
	}
	str.WriteString(fmt.Sprintf(" ID: %v ", m.Header.ID))
	str.WriteString(fmt.Sprintf("Questions: %v", m.Question.Qname))
	str.WriteString("\n")

	return str.String()
}
