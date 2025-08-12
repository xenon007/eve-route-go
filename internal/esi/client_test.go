package esi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStatus(t *testing.T) {
	expectedUA := "test-agent"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/status/" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if ua := r.Header.Get("User-Agent"); ua != expectedUA {
			t.Fatalf("unexpected UA %s", ua)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"players":1,"server_version":"1.0","start_time":"2021-01-01T00:00:00Z","vip":false}`))
	}))
	defer ts.Close()

	c := NewClient(nil, expectedUA)
	c.api.ChangeBasePath(ts.URL)

	st, err := c.Status(context.Background())
	if err != nil {
		t.Fatalf("status error: %v", err)
	}
	expectedTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if st.Players != 1 || st.ServerVersion != "1.0" || !st.StartTime.Equal(expectedTime) || st.Vip {
		t.Fatalf("unexpected status %+v", st)
	}
}
