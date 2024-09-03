// Generated by AI based on https://platform.openai.com/docs/api-reference/chat
package openai

type ChatCompletionRequest struct {
	Messages          []*Message       `json:"messages"`
	Model             string           `json:"model"`
	FrequencyPenalty  float64          `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int64 `json:"logit_bias,omitempty"`
	Logprobs          bool             `json:"logprobs,omitempty"`
	TopLogprobs       int64            `json:"top_logprobs,omitempty"`
	MaxTokens         int64            `json:"max_tokens,omitempty"`
	N                 int64            `json:"n,omitempty"`
	PresencePenalty   float64          `json:"presence_penalty,omitempty"`
	ResponseFormat    *ResponseFormat  `json:"response_format,omitempty"`
	Seed              int64            `json:"seed,omitempty"`
	ServiceTier       string           `json:"service_tier,omitempty"`
	Stop              interface{}      `json:"stop,omitempty"`
	Stream            bool             `json:"stream,omitempty"`
	StreamOptions     *StreamOptions   `json:"stream_options,omitempty"`
	Temperature       float64          `json:"temperature,omitempty"`
	TopP              float64          `json:"top_p,omitempty"`
	Tools             []string         `json:"tools,omitempty"`
	ToolChoice        string           `json:"tool_choice,omitempty"`
	ParallelToolCalls bool             `json:"parallel_tool_calls,omitempty"`
	User              string           `json:"user,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type,omitempty"`
	// NOTE I don't care at the moment
	JsonSchema interface{} `json:"json_schema,omitempty"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type ChatCompletionResponse struct {
	ID                string             `json:"id"`
	Choices           []CompletionChoice `json:"choices"`
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	ServiceTier       string             `json:"service_tier,omitempty"`
	SystemFingerprint string             `json:"system_fingerprint"`
	Object            string             `json:"object"`
	Usage             Usage              `json:"usage"`
}

type CompletionChoice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int64   `json:"index"`
	Message      Message `json:"message"`
	// NOTE I don't care at the moment
	Logprobs interface{} `json:"logprobs,omitempty"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}