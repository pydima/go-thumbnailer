package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestHandleSigTerm(t *testing.T) {
	HandleSigTerm()
	select {
	case <-STOP:
		t.Errorf("Got signal, didn't send it.")
	default:
		break
	}

	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	select {
	case <-STOP:
		break
	case <-time.After(time.Second * 1):
		t.Errorf("Didn't get a signal.")
	}

	// process shouldn't panic if get signal twice
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func TestNotify(t *testing.T) {
	urls := make([]string, 0)
	urls = append(urls, "one")
	urls = append(urls, "two")
	idExpected := "id_string"

	var (
		idResult   string
		urlsResult []string
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ack := Ack{}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Got an error %s", err)
		}

		json.Unmarshal(data, &ack)
		idResult = ack.ID
		urlsResult = ack.Images
	}))
	defer ts.Close()

	ack := NewAck(ts.URL, idExpected, urls)

	if err := Notify(ack); err != nil {
		t.Errorf("got an error %s", err)
	}

	if idResult != idExpected {
		t.Errorf("wrong result, got: %s - expected: %s", idResult, idExpected)
		return
	}

	if len(urls) != len(urlsResult) {
		t.Errorf("wrong result, got: %s - expected: %s", urls, urlsResult)
		return
	}

	for i, _ := range urls {
		if urls[i] != urlsResult[i] {
			t.Errorf("wrong result, got: %s - expected: %s", urls, urlsResult)
			break
		}
	}
}

func TestNotifyMultipleTimes(t *testing.T) {
	urls := make([]string, 0)
	urls = append(urls, "one")

	ack := NewAck("http://localhost:12344/", "id_string", urls)

	if err := Notify(ack, time.Millisecond); err == nil || !strings.Contains(err.Error(), "3 attempts") {
		t.Errorf("Should have got an error with number of attempts to notify client, got %s", err)
	}
}
