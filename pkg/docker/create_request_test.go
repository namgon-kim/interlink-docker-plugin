package docker

import (
	"errors"
	"net/http"
	"testing"
)

func mustHTTPStatus(t *testing.T, err error, want int) {
	t.Helper()
	var hre *httpRequestError
	if !errors.As(err, &hre) {
		t.Fatalf("expected *httpRequestError, got %T: %v", err, err)
	}
	if hre.Status != want {
		t.Fatalf("expected status %d, got %d (msg=%q)", want, hre.Status, hre.Message)
	}
}

func TestParseCreateRequestBody_Object(t *testing.T) {
	body := []byte(`{
		"pod": {
			"metadata": {"uid": "u1", "namespace": "ns", "name": "p1"},
			"spec": {"containers": [], "initContainers": []}
		},
		"container": [],
		"initContainer": [],
		"jobScript": ""
	}`)

	got, err := parseCreateRequestBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if string(got[0].Pod.UID) != "u1" {
		t.Fatalf("expected pod UID u1, got %q", string(got[0].Pod.UID))
	}
}

func TestParseCreateRequestBody_Array(t *testing.T) {
	body := []byte(`[
		{
			"pod": {
				"metadata": {"uid": "u1", "namespace": "ns", "name": "p1"},
				"spec": {"containers": [], "initContainers": []}
			},
			"container": [],
			"initContainer": [],
			"jobScript": ""
		}
	]`)

	got, err := parseCreateRequestBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
}

func TestParseCreateRequestBody_NestedArray(t *testing.T) {
	body := []byte(`[[
		{
			"pod": {
				"metadata": {"uid": "u1", "namespace": "ns", "name": "p1"},
				"spec": {"containers": [], "initContainers": []}
			},
			"container": [],
			"initContainer": [],
			"jobScript": ""
		}
	]]`)

	got, err := parseCreateRequestBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
}

func TestParseCreateRequestBody_EmptyBody(t *testing.T) {
	_, err := parseCreateRequestBody(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	mustHTTPStatus(t, err, http.StatusBadRequest)
}

func TestParseCreateRequestBody_EmptyArray(t *testing.T) {
	_, err := parseCreateRequestBody([]byte(`[]`))
	if err == nil {
		t.Fatal("expected error")
	}
	mustHTTPStatus(t, err, http.StatusUnprocessableEntity)
}

func TestParseCreateRequestBody_InvalidJSONType(t *testing.T) {
	_, err := parseCreateRequestBody([]byte(`123`))
	if err == nil {
		t.Fatal("expected error")
	}
	mustHTTPStatus(t, err, http.StatusBadRequest)
}
