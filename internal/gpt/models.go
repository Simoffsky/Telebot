package gpt

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

func NewRequest(messages []Message) *Request {
	return &Request{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.5,
		Messages:    messages,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Choices []Choice  `json:"choices"`
	Model   string    `json:"model"`
	Error   ErrorResp `json:"error,omitempty"`
}

type ErrorResp struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Param   string `json:"param"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}
