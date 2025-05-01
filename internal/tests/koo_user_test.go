package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/dto"
	"github.com/kootic/koogo/internal/tests/testutils"
)

func TestKooUser(t *testing.T) {
	t.Parallel()

	responses := make([]testutils.TestResponse, 0)

	plan := testutils.TestPlan{
		StepsFactory: []testutils.TestStepFactory{
			func(responses []testutils.TestResponse) testutils.TestStep {
				return testutils.TestStep{
					Name:        "create koo user",
					Path:        "/api/v1/koo/users",
					Method:      http.MethodPost,
					ContentType: "application/json",
					Body: map[string]any{
						"firstName": "John Doe",
					},
					ExpectStatusCode: http.StatusOK,
					ValidateResponse: func(t *testing.T, response testutils.TestResponse) error {
						kooUser, err := testutils.DecodeTestResponse[dto.KooUserResponse](response)
						if err != nil {
							return err
						}

						if kooUser.ID == uuid.Nil {
							return errors.New("user id is nil")
						}

						return nil
					},
				}
			},
			func(responses []testutils.TestResponse) testutils.TestStep {
				newUser, err := testutils.DecodeTestResponse[dto.KooUserResponse](responses[0])
				if err != nil {
					t.Fatalf("failed to get user from previous step")
					return testutils.TestStep{}
				}

				return testutils.TestStep{
					Name:             "get koo user",
					Path:             "/api/v1/koo/users/" + newUser.ID.String(),
					Method:           http.MethodGet,
					ExpectStatusCode: http.StatusOK,
					ValidateResponse: func(t *testing.T, response testutils.TestResponse) error {
						kooUser, err := testutils.DecodeTestResponse[dto.KooUserResponse](response)
						if err != nil {
							return err
						}

						if kooUser.ID != newUser.ID {
							return errors.New("user id does not match")
						}

						if kooUser.FirstName != newUser.FirstName {
							return errors.New("user first name does not match")
						}

						return nil
					},
				}
			},
		},
	}

	testutils.RunTestPlan(t, plan, responses)
}
