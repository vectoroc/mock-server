package model

type HttpTemplate struct {
	Delay        Delay  `json:"delay,omitempty"`
	TemplateType string `json:"templateType,omitempty"`
	Template     string `json:"template,omitempty"`
}
