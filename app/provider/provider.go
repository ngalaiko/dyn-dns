package provider

// Provider allows to manage DNS records.
type Provider interface {
	// Creates a new record.
	Create(*Record) error
	// Updates existing record.
	Update(*Record) error
	// Returns list of existing records.
	Get() ([]*Record, error)
}
