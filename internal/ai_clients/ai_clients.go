package ai_clients

type AIClient interface {
	ChatCompletion([]Message) (Message, int, error)
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}
