package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rehydrate1/shorty/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockURLGetter struct {
	mock.Mock
}

func (m *MockURLGetter) GetURL(ctx context.Context, alias string) (string, error) {
	args := m.Called(ctx, alias)
	return args.String(0), args.Error(1)
}

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		shortKey  string
		mockURL   string
		mockError error
		wantCode  int
		wantURL   string
	}{
		{
			name:      "Success",
			shortKey:  "testKey",
			mockURL:   "https://google.com",
			mockError: nil,
			wantCode:  http.StatusFound, // 302
			wantURL:   "https://google.com",
		},
		{
			name:      "Not Found",
			shortKey:  "unknown",
			mockURL:   "",
			mockError: storage.ErrURLNotFound,
			wantCode:  http.StatusNotFound, // 404
			wantURL:   "",
		},
		{
			name:      "Internal Error",
			shortKey:  "fail",
			mockURL:   "",
			mockError: errors.New("db fall"),
			wantCode:  http.StatusInternalServerError, // 500
			wantURL:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := new(MockURLGetter)

			urlGetterMock.On("GetURL", mock.Anything, tc.shortKey).
				Return(tc.mockURL, tc.mockError)

			gin.SetMode(gin.TestMode)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			
			handler := New(logger, urlGetterMock)

			r := gin.New()
			r.GET("/:shortKey", handler)

			req, _ := http.NewRequest(http.MethodGet, "/"+tc.shortKey, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantCode, w.Code)

			if tc.wantCode == http.StatusFound {
				assert.Equal(t, tc.wantURL, w.Header().Get("Location"))
			}
		})
	}
}