package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"cuboid-challenge/app/router"
)

func MockRequest(method, path string, bodyStr *string) *httptest.ResponseRecorder {
	r := router.Setup()
	w := httptest.NewRecorder()

	var body io.Reader
	if bodyStr != nil {
		body = strings.NewReader(*bodyStr)
	} else {
		body = nil
	}

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, method, path, body)
	r.ServeHTTP(w, req)

	return w
}

func Serialize(d interface{}) (map[string]interface{}, error) {
	out, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize. %w", err)
	}

	var m map[string]interface{}
	if err = json.Unmarshal(out, &m); err != nil {
		return nil, fmt.Errorf("failed to serialize. %w", err)
	}

	return m, nil
}

func SerializeToString(d interface{}) (string, error) {
	out, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("failed to serialize. %w", err)
	}

	return string(out), nil
}

func Deserialize(d string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(d), &m); err != nil {
		return nil, fmt.Errorf("failed to deserialize. %w", err)
	}

	return m, nil
}

func DeserializeList(d string) ([]map[string]interface{}, error) {
	var l []map[string]interface{}
	if err := json.Unmarshal([]byte(d), &l); err != nil {
		return nil, fmt.Errorf("failed to deserialize a list. %w", err)
	}

	return l, nil
}
