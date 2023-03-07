package openai_test

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/arferreira/mentan-blackbox/openai"
)


func TestSendPrompt(t *testing.T) {


    err := godotenv.Load()

    if err != nil {
        return
    }

    tests := []struct {
        name       string
        prompt     string
        wantResult string
        wantErr    bool
    }{
        {
            name:       "Test with valid prompt",
            prompt:     "Hello!",
            wantResult: "{\"document_text\": \"Some response message\"}",
            wantErr:    false,
        },
        {
            name:       "Test with invalid prompt",
            prompt:     "",
            wantResult: "",
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotResult, err := openai.SendPrompt(tt.prompt)
            if (err != nil) != tt.wantErr {
                t.Errorf("SendPrompt() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            assert.Equal(t, tt.wantResult, gotResult)
        })
    }
}
