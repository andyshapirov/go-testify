package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenRequestIsCorrect(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "?count=1&city=moscow", http.NoBody)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// require проставлен в третьем тесте - непосредственно перед манипуляциями с body
	// здесь же NoEmpty требуется по заданию
	assert.NotEmpty(t, responseRecorder.Body)
}

func TestMainHandlerWhenCityIsWrong(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "?count=1&city=dzhambul", http.NoBody)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body)
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	req, err := http.NewRequest(http.MethodGet, "?count=99999999999999&city=moscow", http.NoBody)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NotEmpty(t, responseRecorder.Body)

	cityList := strings.Split(responseRecorder.Body.String(), ",")
	respCount := len(cityList)
	assert.Equal(t, totalCount, respCount)
}