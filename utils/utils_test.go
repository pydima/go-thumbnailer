package utils

import (
	"os"
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
