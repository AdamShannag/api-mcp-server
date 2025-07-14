package types

type Tool struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Args        []Arg   `json:"args"`
	Request     Request `json:"request"`
}

type Arg struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Required     bool   `json:"required"`
	DefaultValue any    `json:"defaultValue"`
	Type         string `json:"type"`
}

type Request struct {
	Host        string            `json:"host"`
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	Secure      bool              `json:"secure"`
	Headers     map[string]string `json:"headers"`
	PathParams  []string          `json:"pathParams"`
	QueryParams []string          `json:"queryParams"`
	Body        string            `json:"body,omitempty"`
}

type Response struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}
