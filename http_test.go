package main

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	_, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegCheckMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:2803/reg-email", nil)
	if err != nil {
		t.Fatal(err)
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		} else {
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Fatal("Method shouldn't be allowed")
			}
			defer resp.Body.Close()
		}
	}
}

func TestRegCheckBadJson(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", nil)
	if err != nil {
		t.Fatal(err)
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		} else {
			if resp.StatusCode != http.StatusBadRequest {
				t.Fatal("Should fails as bad request")
			}
			defer resp.Body.Close()
		}
	}
}
