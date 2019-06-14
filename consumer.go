package main

import (
	"fmt"
	"net"
	"sq"
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
	cfg.TlsV1 = true
	cfg.Set("tls_cert", "/home/zhanglei/cert/server.crt")
	cfg.Set("tls_key", "/home/zhanglei/cert/server.key")

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
	//cfg.TlsV1 = true
	//cfg.Set("tls_cert", "/home/zhanglei/cert/server.csr")
	//cfg.Set("tls_key", "/home/zhanglei/cert/server.key")
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

func main() {
	consumer, err := NewReceiver("202.173.9.22:4150", "tls_test", "tls_test_chan1", "192.168.79.12")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	recvBuff := consumer.GetDataChan()

	var sum int64
	var begin time.Time
	for {
		select {
		case data := <-recvBuff:
			if begin.IsZero() {
				begin = time.Now()
			}
			if string(data) == "11111111111111111111111111111111111111112222222222333333333344444444445555555555666666666677777777771111111111111111111111111111111111111111222222222233333333334444444444555555555566666666667777777777111111111111111111111111111111111111111122222222223333333333444444444455555555556666666666777777777711111111111111111111111111111111111111112222222222333333333344444444445555555555666666666677777777771111111111111111111111111111111111111111222222222233333333334444444444555555555566666666667777777777111111111111111111111111111111111111111122222222223333333333444444444455555555556666666666777777777711111111111111111111111111111111111111112222222222333333333344444444445555555555666666666677777777771111111111111111111111111111111111111111222222222233333333334444444444555555555566666666667777777777111111111111111111111111111111111111111122222222223333333333444444444455555555556666666666777777777711111111111111111111111111111111111111112222222222333333333344444444445555555555666666666677777777771111111111111111111111111111111111111111222222222233333333334444444444555555555566666666667777777777" {
				sum++
			}
			fmt.Println(string(data))
		}

		if sum/100000 != 0 && sum%100000 == 0 {
			fmt.Println("100000 times, duration:", time.Now().Sub(begin))
			begin = time.Now()
		}
	}

}
