package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var random *os.File

// STOP is a channel which is closed when we need to
// stop process, for example when we get a signal.
var STOP chan struct{}

func init() {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal(err)
	}
	random = f
	STOP = make(chan struct{})
}

// Ack contains all information for notifying client
// that image processing task is finished
type Ack struct {
	url    string
	ID     string
	Images []string
}

// NewAck returns pointer to a new initialized Ack structure
func NewAck(url, id string, images []string) *Ack {
	return &Ack{url: url, ID: id, Images: images}
}

// UUID generates new ID for the task, it is used internally
func UUID() string {
	b := make([]byte, 16)
	random.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// DownloadImage donwloads images from the url, and return downloaded image and error
func DownloadImage(u string) ([]byte, error) {
	var data []byte
	u = strings.TrimSpace(u)
	resp, err := http.Get(u)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return data, fmt.Errorf("error status code %d, url: `%s`", resp.StatusCode, u)
	}
	return ioutil.ReadAll(resp.Body)
}

// Notify is using Ack structure to notify client that images is processed
// and provide information about image location
func Notify(ack *Ack, m ...time.Duration) (err error) {
	if ack.url == "" {
		return
	}

	var d time.Duration
	if len(m) > 0 {
		d = m[0]
	} else {
		d = time.Second
	}

	data, err := json.Marshal(ack)
	if err != nil {
		return
	}

	delay := d * 0
	attempts := 3
	for x := attempts; x > 0; x-- {
		<-time.After(delay)
		_, err = http.Post(ack.url, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Printf("got error %s when notified client", err)
			delay = delay + d*1
			continue
		}
		return
	}
	return fmt.Errorf("couldn't notify client, made %d attempts, got: %s", attempts, err)
}

// HandleSignals receives OS signals and notifies goroutines
// that they have to stop processing new tasks. To do this
// functions just closes channel.
func HandleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		close(STOP)
		time.Sleep(time.Second * 15)
		os.Exit(1)
	}()
}
