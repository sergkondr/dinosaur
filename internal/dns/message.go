package dns

// Message format is described in RFC-1035:
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1
type Message struct {
	Header   MessageHeader
	Question MessageQuestion // according to the RFC it is usually equals to 1. I'll fix it someday maybe

	Answer     []MessageResourceRecord
	Authority  []MessageResourceRecord
	Additional []MessageResourceRecord
}

// MessageHeader format is described here:
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type MessageHeader struct {
	ID              uint16
	Flags           Flags
	QueriesCount    uint16
	AnswerCount     uint16
	NameServerCount uint16
	AdditionalCount uint16
}

type Flags struct {
	Query               uint8
	Opcode              uint8
	AuthoritativeAnswer uint8
	Truncation          uint8
	RecursionDesired    uint8
	RecursionAvailable  uint8
	Z                   uint8
	ResponseCode        uint8
}

// MessageQuestion format is described here:
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
type MessageQuestion struct {
	Qname  []string
	Qtype  uint16
	Qclass uint16
}

// MessageResourceRecord format is described here:
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3
// The answer, authority, and additional sections all share the same format
type MessageResourceRecord struct {
	Name  []string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []int
}
