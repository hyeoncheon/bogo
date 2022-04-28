package common

const (
	// UNKNOWN is a kind of "I don't know" from the perspective of a specific
	// driver. For example, GCE meta client can answer with UNKNOWN if they
	// could not detect as they are on GCE.
	UNKNOWN = "Unknown"
	// NOWHERE has different meaning. For example, if this keyword is used,
	// it should indicate that the environment is explicitly not supported.
	NOWHERE = "Nowhere"
)

// MetaClient is an interface to be implemented for each cloud provider.
type MetaClient interface {
	// WhereAmI returns the name of CSP.
	WhereAmI() string
	// InstanceName returns the current VM's instance name string.
	InstanceName() string
	// ExternalIP returns the instance's primary external (public) IP address.
	ExternalIP() string
	// Zone returns the current VM's zone, such as "asia-northeast3-a".
	Zone() string
	// AttributeValues gets comma or space separated metadata values
	// It is a wrapper of AttributeValue(string)
	AttributeValues(key string) []string
	// AttributeValue gets bare string for the given metadata key. It will
	// check the metadata for the instance first and will return the value
	// if the metadata with the same key exists, otherwise it will check
	// the project's metadata with the same logic.
	AttributeValue(key string) string
	// AttributeSSV gets space separated values for the given metadata key
	// It is a wrapper of AttributeValue(string)
	AttributeSSV(key string) []string
	// AttributeCSV gets comma separated values for the given metadata key
	// It is a wrapper of AttributeValue(string)
	AttributeCSV(key string) []string
}
