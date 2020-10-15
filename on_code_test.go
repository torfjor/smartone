package smartone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCodeHandler_Invalid_Method(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	want := http.StatusMethodNotAllowed
	wantHeader := "POST"

	OnCodeHandler(func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error) {
		return OnCodeResponse{}, nil
	})(res, req)

	got := res.Code
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}

	gotHeader := res.Header().Get("Allow")
	if !strings.Contains(gotHeader, wantHeader) {
		t.Errorf("got allow header %q, want %q", gotHeader, wantHeader)
	}
}

func TestOnCodeHandler_Invalid_JSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"invalid"}`)))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	wantCode := http.StatusBadRequest
	OnCodeHandler(func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error) {
		return OnCodeResponse{}, nil
	})(res, req)

	if res.Code != wantCode {
		t.Errorf("got status %d, want %d", res.Code, wantCode)
	}
}

func TestOnCodeHandler_Invalid_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "text/html")
	res := httptest.NewRecorder()

	OnCodeHandler(func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error) {
		return OnCodeResponse{}, nil
	})(res, req)
	wantCode := http.StatusBadRequest

	if res.Code != wantCode {
		t.Errorf("got code %d, want %d", res.Code, wantCode)
	}
}

func TestOnCodeHandler_Valid_Request(t *testing.T) {
	var body = `{
  "Code": "1234",
  "DoorId": 1,
  "SLXId": "445e831c-abd3-4178-9005-112233445566",
  "SERNO": "112233445566"
}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	OnCodeHandler(func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error) {
		if req.Code != "1234" {
			t.Errorf("got code %q, want %q", req.Code, "1234")
		}
		if req.DoorID != 1 {
			t.Errorf("got doorID %d, want %d", req.DoorID, 1)
		}
		if req.SLXID != "445e831c-abd3-4178-9005-112233445566" {
			t.Errorf("got SLXID %q, want %q", req.SLXID, "445e831c-abd3-4178-9005-112233445566")
		}
		if req.SerialNumber != "112233445566" {
			t.Errorf("got SerialNumber %q, want %q", req.SerialNumber, "112233445566")
		}

		return OnCodeResponse{Result: true}, nil
	})(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", res.Code, http.StatusOK)
	}

	var response OnCodeResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("err=%v", err)
	}

	if response.Result != true {
		t.Errorf("got %v, want %v", response.Result, true)
	}
}

func TestOnCodeHandler_Valid_Request_Err_OnCodeFunc(t *testing.T) {
	var body = `{
  "Code": "1234",
  "DoorId": 1,
  "SLXId": "445e831c-abd3-4178-9005-112233445566",
  "SERNO": "112233445566"
}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	OnCodeHandler(func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error) {
		return OnCodeResponse{}, fmt.Errorf("some error")
	})(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", res.Code, http.StatusInternalServerError)
	}
}
