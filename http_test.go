package main

import (
	"net/http"
	"testing"
	"bytes"
)

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


func TestRegCheckBadEmail(t *testing.T) {
	var jsonStr = []byte(`{"email":";X@@example", "nick":"any"}`)
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", bytes.NewBuffer(jsonStr))
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

func TestRegCheckBadNickName(t *testing.T) {
	var jsonStr = []byte(`{"email":"johndoe@any.lynx", "nick":"john@"}`)
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", bytes.NewBuffer(jsonStr))
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

func TestRegCheckShouldPass(t *testing.T) {
	var jsonStr = []byte(`{"email":"johndoe@any.lynx", "nick":"john"}`)
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		} else {
			if resp.StatusCode != http.StatusOK {
				t.Fatal("Should pass:")
			}
			defer resp.Body.Close()
		}
	}
}

func TestRegCheckEmailConflict(t *testing.T) {
	var jsonStr = []byte(`{"email":"johndoe@any.lynx", "nick":"any"}`)
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		} else {
			if resp.StatusCode != http.StatusConflict {
				t.Fatal("Should fails email conflict")
			}
			defer resp.Body.Close()
		}
	}
}

func TestRegCheckNickConflict(t *testing.T) {
	var jsonStr = []byte(`{"email":"any-johndoe@any.lynx", "nick":"john"}`)
	req, err := http.NewRequest("POST", "http://localhost:2803/reg-email", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		} else {
			if resp.StatusCode != http.StatusConflict {
				t.Fatal("Should fails email conflict")
			}
			defer resp.Body.Close()
		}
	}
}