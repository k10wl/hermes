package ai_clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/k10wl/hermes/internal/ai_clients/openai"
)

func TestOpenAICompletion(t *testing.T) {
	inputMessages := []*Message{{Role: "user", Content: "stuff"}}
	outputMessages := []*Message{{Role: "assistant", Content: "response"}}
	temperature := float64(1)
	maxTokens := int64(1000)

	var lastRequest *http.Request

	type input struct {
		messages   []*Message
		parameters Parameters
		getter     getter
	}
	type testCase struct {
		name             string
		input            input
		apiKey           string
		expectedResponse AIResponse
		expectedHeaders  map[string]string
		expectErr        bool
	}

	table := []testCase{
		{
			name: "should call complete with given params",
			input: input{
				messages: inputMessages,
				parameters: Parameters{
					Temperature: &temperature,
					MaxTokens:   &maxTokens,
					Model:       "gpt-4o-mini",
				},
				getter: func(
					url string,
					body io.Reader,
					fillHeaders func(*http.Request) error,
				) ([]byte, error) {
					lastRequest, _ = http.NewRequest(http.MethodPost, "", nil)
					err := fillHeaders(lastRequest)
					if err != nil {
						return nil, err
					}
					res := openai.ChatCompletionResponse{
						ID:                "1",
						SystemFingerprint: "SystemFingerprint",
						Created:           1234,
						Model:             "gpt-4o-mini",
						Object:            "Object",
						Usage: openai.Usage{
							PromptTokens:     200,
							CompletionTokens: 20,
							TotalTokens:      40,
						},
						Choices: []openai.CompletionChoice{{
							FinishReason: "stop",
							Index:        0,
							Message: openai.Message{
								Role:    "assistant",
								Content: "response",
							},
						}},
					}
					marshaled, err := json.Marshal(res)
					if err != nil {
						panic("bad test setup, error in marshaling response")
					}
					return marshaled, nil
				},
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer SECRET",
			},
			expectedResponse: AIResponse{
				Messages: outputMessages,
				TokensUsage: TokensUsage{
					Input:  200,
					Output: 20,
				},
			},
			apiKey: "SECRET",
		},

		{
			name: "should error if api key was not provided",
			input: input{
				messages: inputMessages,
				parameters: Parameters{
					Temperature: &temperature,
					MaxTokens:   &maxTokens,
					Model:       "gpt-4o-mini",
				},
				getter: func(
					url string,
					body io.Reader,
					fillHeaders func(*http.Request) error,
				) ([]byte, error) {
					lastRequest, _ = http.NewRequest(http.MethodPost, "", nil)
					err := fillHeaders(lastRequest)
					if err != nil {
						return nil, err
					}
					res := openai.ChatCompletionResponse{
						ID:                "1",
						SystemFingerprint: "SystemFingerprint",
						Created:           1234,
						Model:             "gpt-4o-mini",
						Object:            "Object",
						Usage: openai.Usage{
							PromptTokens:     200,
							CompletionTokens: 20,
							TotalTokens:      40,
						},
						Choices: []openai.CompletionChoice{{
							FinishReason: "stop",
							Index:        0,
							Message: openai.Message{
								Role:    "assistant",
								Content: "response",
							},
						}},
					}
					marshaled, err := json.Marshal(res)
					if err != nil {
						panic("bad test setup, error in marshaling response")
					}
					return marshaled, nil
				},
			},
			expectErr: true,
		},

		{
			name: "should return error if getter errored",
			input: input{
				messages: inputMessages,
				parameters: Parameters{
					Temperature: &temperature,
					MaxTokens:   &maxTokens,
					Model:       "gpt-4o-mini",
				},
				getter: func(
					url string,
					body io.Reader,
					fillHeaders func(*http.Request) error,
				) ([]byte, error) {
					return []byte{}, fmt.Errorf("it just errored")
				},
			},
			expectErr: true,
		},
	}

out:
	for _, test := range table {
		actual, err := newClientOpenAI(test.apiKey).complete(
			test.input.messages,
			test.input.parameters,
			test.input.getter,
		)
		if test.expectErr {
			if err == nil {
				t.Errorf("%q - expected to error but didn't\n", test.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q - unexpected error: %v\n", test.name, err)
			continue
		}
		for key, val := range test.expectedHeaders {
			if header := lastRequest.Header.Get(key); header != val {
				t.Errorf(
					"%q - bad header for %q.\nexpected: %q\nactual:   %q\n\n",
					test.name,
					key,
					val,
					header,
				)
				continue out
			}
		}
		if actual == nil {
			t.Errorf("%q - actual result is nil\n\n", test.name)
			continue
		}
		if !reflect.DeepEqual(test.expectedResponse, *actual) {
			t.Errorf(
				"%q - bad return from client.\nexpected: %+v\nactual:   %+v\n\n",
				test.name,
				test.expectedResponse,
				*actual,
			)
			continue
		}
	}
}
