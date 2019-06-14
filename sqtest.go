package main

import (
	"flag"
	"fmt"
	"net"
	"sq"
	"strconv"
	"time"
)

const (
	BuffleLen      int = 0
	MaxResendCount     = 3
)
const MaxResendWaitTime = 10 * time.Second

type SqReceiver struct {
	*sq.Consumer
	dataChan chan []byte
}

func NewReceiver(sqdAddr, topic, channel, laddr string) (*SqReceiver, error) {
	cfg := sq.NewConfig()
	cfg.LocalAddr = &net.TCPAddr{IP: net.ParseIP(laddr)}

	consumer, err := sq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	receiver := &SqReceiver{
		Consumer: consumer,
		dataChan: make(chan []byte, BuffleLen),
	}

	consumer.AddHandler(receiver)
	err = consumer.ConnectToSQD(sqdAddr)
	if err != nil {
		consumer.Stop()
		return nil, err
	}

	return receiver, nil
}
func NewTlsReceiver(sqdAddr, topic, channel, laddr string) (*SqReceiver, error) {
	cfg := sq.NewConfig()
	cfg.LocalAddr = &net.TCPAddr{IP: net.ParseIP(laddr)}
	cfg.TlsV1 = true

	consumer, err := sq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	receiver := &SqReceiver{
		Consumer: consumer,
		dataChan: make(chan []byte, BuffleLen),
	}

	consumer.AddHandler(receiver)
	err = consumer.ConnectToSQD(sqdAddr)
	if err != nil {
		consumer.Stop()
		return nil, err
	}

	return receiver, nil
}

func (r *SqReceiver) HandleMessage(m *sq.Message) error {
	r.dataChan <- m.Body
	return nil
}

func (r *SqReceiver) GetDataChan() <-chan []byte {
	return r.dataChan
}

type SqSender struct {
	*sq.Producer
	dataChan chan []byte
	topic    string
}

func NewSender(sqdAddr, topic, laddr string) (*SqSender, error) {
	cfg := sq.NewConfig()
	cfg.LocalAddr = &net.TCPAddr{IP: net.ParseIP(laddr)}
	producer, err := sq.NewProducer(sqdAddr, cfg)
	if err != nil {
		return nil, err
	}

	sender := SqSender{
		Producer: producer,
		dataChan: make(chan []byte, BuffleLen),
		topic:    topic,
	}

	go sender.loop()
	return &sender, nil
}

func NewTlsSender(sqdAddr, topic, laddr string) (*SqSender, error) {
	cfg := sq.NewConfig()
	cfg.LocalAddr = &net.TCPAddr{IP: net.ParseIP(laddr)}
	cfg.TlsV1 = true
	producer, err := sq.NewProducer(sqdAddr, cfg)
	if err != nil {
		return nil, err
	}

	sender := SqSender{
		Producer: producer,
		dataChan: make(chan []byte, BuffleLen),
		topic:    topic,
	}

	go sender.loop()
	return &sender, nil
}

func (s *SqSender) loop() {
	doneCh := make(chan *sq.ProducerTransaction)

	for {
		cmd := <-s.dataChan
		resendWaitTime := time.Second
		for {
			err := s.PublishAsync(s.topic, cmd, doneCh)
			if err == nil {
				tran := <-doneCh
				if tran.Error == nil {
					break
				}
			}

			<-time.After(resendWaitTime)
			if resendWaitTime*2 < MaxResendWaitTime {
				resendWaitTime = 2 * resendWaitTime
			} else {
				resendWaitTime = MaxResendWaitTime
			}
		}
	}
}

func (s *SqSender) GetDataChan() chan<- []byte {
	return s.dataChan
}

var (
	local   string
	sqd     string
	topic   string
	channel string
	msg     string
	timeout int
)

func init() {
	flag.StringVar(&local, "local", "", "local ip")
	flag.StringVar(&sqd, "sqd", "", "sqd address")
	flag.StringVar(&topic, "topic", "", "topic name")
	flag.StringVar(&channel, "channel", "", "channel name")
	flag.StringVar(&msg, "msg", "", "msg will send")
	flag.IntVar(&timeout, "timeout", 0, "timeout for receive")
}

func main() {
	flag.Parse()
	if local == "" ||
		sqd == "" ||
		topic == "" {
		fmt.Println("please check param...")
		return
	}
	if channel != "" {
		re, err := NewReceiver(sqd, topic, channel, local)
		if err != nil {
			panic(err.Error())
		}
		for {
			ch := re.GetDataChan()
			data := <-ch
			fmt.Println(string(data))
			<-time.After(time.Duration(timeout) * time.Second)
		}
	} else {
		if msg == "" {
			fmt.Println("message is nil")
			return
		}
		re, _ := NewSender(sqd, topic, local)
		ch := re.GetDataChan()
		for i := 1; i <= 100; i++ {
			ch <- []byte(msg + strconv.Itoa(i))
		}
		<-time.After(1 * time.Second)
	}

}
