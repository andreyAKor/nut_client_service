package get

// UPS contains information about a specific UPS provided by the NUT instance.
type UPS struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Master         bool       `json:"master"`
	NumberOfLogins int        `json:"numberOfLogins"`
	Clients        []string   `json:"clients"`
	Variables      []Variable `json:"variables"`
	Commands       []Command  `json:"commands"`
}

// Variable describes a single variable related to a UPS.
type Variable struct {
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	Type          string      `json:"type"`
	Description   string      `json:"description"`
	Writeable     bool        `json:"writeable"`
	MaximumLength int         `json:"maximumLength"`
	OriginalType  string      `json:"originalType"`
}

// Command describes an available command for a UPS.
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
