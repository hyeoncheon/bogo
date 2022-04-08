package common

const (
	NOWHERE = "Nowhere"
	UNKNOWN = "Unknown"
)

// MetaClient is an interface to be implemented for each cloud provider.
type MetaClient interface {
	WhereAmI() string
	InstanceName() string
	ExternalIP() string
	Zone() string
	// AttributeValues gets comma or space separated metadata values
	AttributeValues(string) []string
	// AttributeValue gets bare string for the given metadata key
	AttributeValue(string) string
	// AttributeSSV gets space separated values for the given metadata key
	AttributeSSV(string) []string
	// AttributeCSV gets comma separated values for the given metadata key
	AttributeCSV(string) []string
}
