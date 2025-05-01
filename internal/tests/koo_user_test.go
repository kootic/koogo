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

	plan := testutils.TestPlan{
		func(globalVars map[string]any) testutils.TestStep {
			return testutils.TestStep{
				Name:        "create koo user",
				Path:        "/api/v1/koo/users",
				Method:      http.MethodPost,
				ContentType: "application/json",
				Body: map[string]any{
					"firstName": "John Doe",
				},
				ExpectStatusCode: http.StatusOK,
				ValidateResponse: func(t *testing.T, response testutils.TestResponse, globalVars map[string]any) error {
					kooUser, err := testutils.DecodeTestResponse[dto.KooUserResponse](response)
					if err != nil {
						return err
					}

					if kooUser.ID == uuid.Nil {
						return errors.New("user id is nil")
					}

					globalVars["newUser"] = kooUser

					return nil
				},
			}
		},
		func(globalVars map[string]any) testutils.TestStep {
			newUser := globalVars["newUser"].(dto.KooUserResponse)

			return testutils.TestStep{
				Name:             "get koo user",
				Path:             "/api/v1/koo/users/" + newUser.ID.String(),
				Method:           http.MethodGet,
				ExpectStatusCode: http.StatusOK,
				ValidateResponse: func(t *testing.T, response testutils.TestResponse, globalVars map[string]any) error {
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
	}

	testutils.RunTestPlan(t, plan)
}
