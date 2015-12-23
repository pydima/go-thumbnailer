package utils

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestHandleSigTerm(t *testing.T) {
	done := HandleSigTerm()
	select {
	case <-done:
		t.Errorf("Got signal, didn't send it.")
	default:
		break
	}

	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	select {
	case <-done:
		break
	case <-time.After(time.Second * 1):
		t.Errorf("Didn't get a signal.")
	}

	// process shouldn't panic if get signal twice
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
