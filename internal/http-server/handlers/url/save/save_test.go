package save

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLSaver struct {
	mock.Mock
}

func (m *MockURLSaver) SaveURL(ctx context.Context, alias, urlToSave string) error {
	args := m.Called(ctx, alias, urlToSave)
	return args.Error(0)
}

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name string
		inputURL string
		mockError error
		wantCode int
		wantBody string
	}{
		{
			name:      "Success",
			inputURL:  "https://google.com",
			mockError: nil,
			wantCode:  http.StatusCreated,
		},
		{
			name:      "Empty URL",
			inputURL:  "",
			mockError: nil,
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "Save Error",
			inputURL:  "https://google.com",
			mockError: errors.New("db error"),
			wantCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := new(MockURLSaver)

			if tc.inputURL != "" {
				urlSaverMock.On("SaveURL", mock.Anything, mock.AnythingOfType("string"), tc.inputURL).
					Return(tc.mockError)
			}

			gin.SetMode(gin.TestMode)
			
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			handler := New(logger, urlSaverMock, "http://localhost:8080")

			r := gin.New()
			r.POST("/shorten", handler)

			bodyData := map[string]string{"url": tc.inputURL}
			body, _ := json.Marshal(bodyData)

			req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}