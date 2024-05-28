package ai_clients

type Message struct {
	Content string
	Role    string
}

type AIClient interface {
	/*
		Returns:
		 - result message
		 - amount of messages used for the completion due to token limit
		 - error
	*/
	ChatCompletion([]Message) (Message, int, error)
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}
