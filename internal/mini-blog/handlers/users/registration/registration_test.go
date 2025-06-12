package registration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"mini-blog/internal/mini-blog/handlers/users/registration"
	"mini-blog/internal/mini-blog/handlers/users/registration/mocks"
)

func TestRegistrationHandler(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			username: "user1",
		},
		{
			name:      "Empty Username",
			respError: "field 'Username' is required",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			userSaverMock := mocks.NewUserSaver(t)

			if c.respError == "" || c.mockError != nil {
				userSaverMock.On("SaveUser", mock.AnythingOfType("context.backgroundCtx"), c.username).
					Return(int64(1), c.mockError).
					Once()
			}

			handler := registration.New(userSaverMock)

			input := fmt.Sprintf(`{"username": "%s"}`, c.username)

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp registration.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, c.respError, resp.Error)
		})
	}
}
