package ai_clients

type Message struct {
	Content string
	Role    string
}

type AIClient interface {
	ChatCompletion([]Message) (Message, error)
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}
