package provider

// RecordType is a type of a DNS record.
type RecordType uint

// DNS record types.
const (
	RecordTypeA RecordType = iota
)

// Record is a single DNS record.
type Record struct {
	Type  RecordType
	Name  string
	Value string
}
