package model

// Body - request body matcher.
type Body struct {
	Not         bool   `json:"not,omitempty"`
	SubString   bool   `json:"subString,omitempty"`
	Type        string `json:"type,omitempty"`
	Base64Bytes string `json:"base64Bytes,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Json        string `json:"json,omitempty"`
	MatchType   string `json:"matchType,omitempty"`
	JsonSchema  string `json:"jsonSchema,omitempty"`
	JsonPath    string `json:"jsonPath,omitempty"`
	//Parameters  *KeyToMultiValue `json:"parameters,omitempty"`
	Regex     string `json:"regex,omitempty"`
	String    string `json:"string,omitempty"`
	Xml       string `json:"xml,omitempty"`
	XmlSchema string `json:"xmlSchema,omitempty"`
	Xpath     string `json:"xpath,omitempty"`
}
