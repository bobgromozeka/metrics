package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bobgromozeka/metrics/internal/server/storage"
)

const Key = "key"

func TestUpdateJSON_BadRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", nil)
	httpW := httptest.NewRecorder()

	stor := storage.NewMemory()
	server := New(
		stor, Config{
			HashKey: Key,
		}, []byte{},
	)

	server.ServeHTTP(httpW, req)
	result := httpW.Result()
	defer result.Body.Close()

	respBody, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(t, "Bad request: EOF\n", string(respBody))
	assert.Equal(t, http.StatusBadRequest, httpW.Code)
}

func TestUpdateJSON_WrongMetricsType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id", "type":"random"}`))
	httpW := httptest.NewRecorder()
	defer httpW.Result().Body.Close()

	stor := storage.NewMemory()
	server := New(
		stor, Config{
			HashKey: Key,
		}, []byte{},
	)

	server.ServeHTTP(httpW, req)

	result := httpW.Result()
	defer result.Body.Close()

	respBody, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(t, "Wrong metrics type\n", string(respBody))
	assert.Equal(t, http.StatusBadRequest, httpW.Code)
}

func TestUpdateJSON_CounterType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id","type":"counter","delta":22}`))
	httpW := httptest.NewRecorder()
	defer httpW.Result().Body.Close()

	stor := storage.NewMemory()
	server := New(
		stor, Config{
			HashKey: Key,
		}, []byte{},
	)

	stor.AddCounter(context.Background(), "id", 20)

	server.ServeHTTP(httpW, req)

	result := httpW.Result()
	defer result.Body.Close()

	respBody, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(t, "application/json", httpW.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":"id","type":"counter","delta":42}`+"\n", string(respBody))
	assert.Equal(t, http.StatusOK, httpW.Code)
}

func TestUpdateJSON_GaugeType(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id": "id","type":"gauge","value":33}`))
	httpW := httptest.NewRecorder()
	defer httpW.Result().Body.Close()

	stor := storage.NewMemory()
	server := New(
		stor, Config{
			HashKey: Key,
		}, []byte{},
	)

	stor.SetGauge(context.Background(), "id", 123.123)

	server.ServeHTTP(httpW, req)

	result := httpW.Result()
	defer result.Body.Close()

	respBody, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(t, "application/json", httpW.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":"id","type":"gauge","value":33}`+"\n", string(respBody))
	assert.Equal(t, http.StatusOK, httpW.Code)
}
