package openai_test

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/arferreira/mentan-blackbox/openai"
)

type EbookInfoProduct struct {
    Title         string `json:"title"`
    Niche         string `json:"niche"`
    Organization string `json:"organizationId"`
    Product       string `json:"productId"`
    Format        string `json:"format"`
}

func TestSendPrompt(t *testing.T) {
    ebook := EbookInfoProduct{
        Title:         "A e-book about Golang programming language",
        Niche: "programming",
        Organization: "123445678",
        Product:       "12345678",
        Format:        "ebook",
    }

    prompt := fmt.Sprintf("Create a ebook for me about %s using this %s", ebook.Niche, ebook.Title)
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
            prompt:     prompt,
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
