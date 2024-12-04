package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	// оставил только параметрическую чать url, т.к. хендлер mainHandle не чувствителен к пути на ресурс
	req, err := http.NewRequest(http.MethodGet, "?count=999999&city=moscow", http.NoBody)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// здесь нужно добавить необходимые проверки
	require.NotEmpty(t, responseRecorder.Body)

	if http.StatusOK == responseRecorder.Code {
		if reqCount, _ := strconv.Atoi(req.URL.Query().Get("count")); reqCount >= totalCount {
			cityList := strings.Split(responseRecorder.Body.String(), ",")
			respCount := len(cityList)
			assert.Equal(t, totalCount, respCount)
		}
	}

	if _, ok := cafeList[req.URL.Query().Get("city")]; !ok {
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Equal(t, "wrong city value", responseRecorder.Body.String())
	}
}
