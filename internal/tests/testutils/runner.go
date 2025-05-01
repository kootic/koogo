package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type TestPlan []TestStepFactory

// TestStepFactory is a function that receives the responses from previous steps
// and returns a new testStep in order to build steps that depend on previous steps.
//
// For example, if the first step creates a user, and the second step gets the user,
// the second step's testStepFactory will a slice of testResponse with length 1,
// where the first element is the response from the first step.
type TestStepFactory func(globalVars map[string]any) TestStep

// TestResponse is a struct that contains the raw response body from a test step
// and anything else that might be needed to build the following steps can be added here.
type TestResponse struct {
	RawBody []byte
}

// DecodeTestResponse decodes the raw response body into a struct.
func DecodeTestResponse[T any](r TestResponse) (T, error) {
	var decoded T

	err := json.Unmarshal(r.RawBody, &decoded)
	if err != nil {
		return decoded, err
	}

	return decoded, nil
}

type TestStep struct {
	Name        string
	Path        string
	Method      string
	ContentType string
	// Body should be a map[string]any most of the time so we can also validate the JSON marshalling
	Body             any
	ExpectStatusCode int
	ValidateResponse func(t *testing.T, response TestResponse, globalVars map[string]any) error
}

func RunTestPlan(t *testing.T, plan TestPlan) {
	t.Helper()

	globalVars := make(map[string]any)
	for _, stepFactory := range plan {
		runTestStep(t, stepFactory, globalVars)
	}
}

func runTestStep(t *testing.T, stepFactory TestStepFactory, globalVars map[string]any) {
	t.Helper()

	step := stepFactory(globalVars)

	var newResponse TestResponse

	t.Run(step.Name, func(t *testing.T) {
		// Marshal the body if it exists
		var body io.Reader
		if step.Body != nil {
			jsonBody, err := json.Marshal(step.Body)
			if err != nil {
				t.Fatalf("Failed to marshal body: %v", err)
			}

			body = bytes.NewReader(jsonBody)
		}

		req, _ := http.NewRequest(step.Method, step.Path, body)
		if step.ContentType != "" {
			req.Header.Set("Content-Type", step.ContentType)
		}

		// Send the request and get the response
		resp, err := TestApp.FiberApp().Test(req)
		if err != nil {
			t.Fatalf("Failed to send test request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		// Read the response body
		newResponse, err = getTestResponse(resp)
		if err != nil {
			t.Fatalf("Failed to get test response: %v", err)
		}

		// Run validations
		if step.ExpectStatusCode != 0 {
			if resp.StatusCode != step.ExpectStatusCode {
				t.Fatalf("Expected status %d, got %d", step.ExpectStatusCode, resp.StatusCode)
			}
		}

		if step.ValidateResponse != nil {
			err := step.ValidateResponse(t, newResponse, globalVars)
			if err != nil {
				t.Fatalf("Failed to validate response: %v", err)
			}
		}
	})
}

// getTestResponse reads the response body and returns a testResponse.
func getTestResponse(resp *http.Response) (TestResponse, error) {
	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResponse{}, err
	}

	return TestResponse{
		RawBody: respBody,
	}, nil
}
