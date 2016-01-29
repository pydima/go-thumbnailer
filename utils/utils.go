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
	"syscall"
	"time"
)

var Random *os.File
var STOP chan struct{}

func init() {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal(err)
	}
	Random = f
	STOP = make(chan struct{})
}

type Ack struct {
	url    string `json:"-"`
	ID     string
	Images []string
}

func NewAck(url, id string, images []string) *Ack {
	return &Ack{url: url, ID: id, Images: images}
}

func UUID() string {
	b := make([]byte, 16)
	Random.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func DownloadImage(u string) ([]byte, error) {
	var data []byte
	resp, err := http.Get(u)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return data, fmt.Errorf("error status code %d.", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func Notify(ack *Ack) (err error) {
	data, err := json.Marshal(ack)
	if err != nil {
		return
	}
	http.Post(ack.url, "application/json", bytes.NewReader(data))
	return
}

func HandleSigTerm() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		close(STOP)
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}()
}
