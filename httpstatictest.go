package main

import (
	"cement/log"
	"quark"
	"quark/httpcmd"
	"quark/httpstatic"
	"quark/registry"
)

func main() {
	registry, err := registry.NewEtcdRegistry(cfg.IP, cfg.EtcdCluster)
	if err != nil {
		panic("create registry failed:" + err.Error())
	}
	logger := log.NewLog4jLogger(cfg.LogFile, log.Debug, 0, 0)
	s := newImageCreate(logger)
	e := &quark.EndPoint{
		Name: "image_service",
		IP:   cfg.IP,
		Port: cfg.Port,
	}
	if err := httpstatic.Run(s, registry, e); err != nil {
		panic("start image service failed:" + err.Error())
	}
}

type imageCreate struct {
	logger log.Logger
}

func newImageCreate(log log.Logger) *imageCreate {
	return &imageCreate{log}
}

func (p *imageCreate) HandleGetFile(user string, c *httpstatic.GetFile) (*httpstatic.Attachment, error) {

	var content bytes.Buffer
	var id string
	var err error
	opt := c.Params["opt"]
	if opt == "get_image_id" {
		id = captcha.New()
		p.logger.Debug("create image id:%s", id)
		return &httpstatic.Attachment{
			FileName: id + ".txt",
			FileType: "text/txt",
			NoAttach: true,
			Content:  []byte(id),
		}, nil
	} else {
		return nil, errors.New("unknown command.")
	}

	if err != nil {
		p.logger.Error("create image failed::%s", err.Error())
		return nil, errors.New("create image failed.")
	}

	return &httpstatic.Attachment{
		FileName: id + ".png",
		FileType: "image/png",
		NoAttach: true,
		Content:  content.Bytes(),
	}, nil
}
