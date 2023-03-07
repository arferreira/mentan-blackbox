package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arferreira/mentan-blackbox/openai"
	"github.com/stretchr/testify/assert"
)

func TestSendPrompt_Success(t *testing.T) {
    // expected result
    expectedResult := "This is the text response from OpenAI."

    // configure the data to return as part of the API response
    responseJson := fmt.Sprintf(
        `{"choices":[{"text":"%s"}]}`,
        expectedResult,
    )

    // create a test server which returns the desired JSON payload
    srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

        assert.Equal(t, "Bearer <API_KEY_HERE>", req.Header.Get("Authorization"))

        var openAiPrompt openai.OpenAiPrompt
        err := json.NewDecoder(req.Body).Decode(&openAiPrompt)
        assert.Nil(t, err)
        assert.Equal(t, "test-prompt", openAiPrompt.Prompt)
        assert.Equal(t, 0.5, openAiPrompt.Temperature)
        assert.Equal(t, 10, openAiPrompt.MaxTokens)

        assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

        res.Write([]byte(responseJson))
    }))
    defer srv.Close()

    // set the fake API base URL in the library
    os.Setenv("OPENAI_BASE_URL", srv.URL)

    // test calling the function
    result, err := openai.SendPrompt(openai.OpenAiPrompt{
        Prompt:      "test-prompt",
        Temperature: 0.5,
        MaxTokens:   10,
    })

    // check that no error occurred and that the expected result was returned
    assert.Nil(t, err)
    assert.Equal(t, expectedResult, result)
}

func TestSendPrompt_Error(t *testing.T) {
    // configure the server to return an error status code
    srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(500)
    }))
    defer srv.Close()

    // ensure that the test function doesn't have access to a valid API key
    os.Unsetenv("OPENAI_API_KEY")

    // set the fake API base URL in the library
    os.Setenv("OPENAI_BASE_URL", srv.URL)

    // test calling the function
    result, err := openai.SendPrompt(openai.OpenAiPrompt{
        Prompt:      "test-prompt",
        Temperature: 0.5,
        MaxTokens:   10,
    })

    // check that an error occurred and no result was returned
    assert.NotNil(t, err)
    assert.Nil(t, result)
}
