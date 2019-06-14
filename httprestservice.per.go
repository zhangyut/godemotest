package personalservice

import (
	"cement/log"
	"errors"
	"quark"
	"quark/rest"
	"zcloud-go/personalservice/resource"
	"zcloud-go/usermanager"
)

type RestHandle struct {
	logger log.Logger
	proxy  quark.ServiceProxy
}

func NewRestHandle(logger log.Logger, proxy quark.ServiceProxy) *RestHandle {
	return &RestHandle{logger, proxy}
}

func (p *RestHandle) HandleGet(store rest.ResourceStore, t *quark.Task) *quark.TaskResult {
	if len(t.Cmds) > 1 {
		p.logger.Error("not support batch cmds")
		return t.Failed(errors.New("not support batch cmds"))
	}
	tx, _ := store.Begin()
	defer tx.Commit()
	c, _ := t.Cmds[0].(*rest.GetCmd)
	results, _ := tx.Get(rest.ResourceType(c.ResourceType), c.Conds)
	return t.SucceedWithResult(results)
}

func (p *RestHandle) HandlePost(store rest.ResourceStore, t *quark.Task) *quark.TaskResult {
	var err error
	if len(t.Cmds) > 1 {
		p.logger.Error("not support batch cmds")
		return t.Failed(errors.New("not support batch cmds"))
	}
	tmp := quark.NewTask()
	tmp.User = t.User
	tmp.AddCmd(&rest.GetCmd{
		ResourceType: "zdnsuser",
		Conds:        map[string]interface{}{"id": t.User},
	})
	ret := ""
	zdnsuser := []usermanager.Zdnsuser{}
	p.proxy.HandleTask(t, &zdnsuser, &ret)
	if ret != "" {
		p.logger.Error("get user error:" + ret)
		return t.Failed(errors.New("get user error"))
	}
	if len(zdnsuser) != 1 {
		p.logger.Error("get user error.")
		return t.Failed(errors.New("get user error"))
	}
	tx, _ := store.Begin()

	c, _ := t.Cmds[0].(*rest.PostCmd)
	personalData, _ := c.NewResource.(*resource.PersonalData)
	personalData.Id = personalData.Zdnsuser
	personalData.Username = zdnsuser[0].Name
	newResource, err := tx.Insert(c.NewResource)
	if err != nil {
		tx.RollBack()
		p.logger.Error("add persional data failed:" + err.Error())
		return t.Failed(errors.New("add persional data failed"))
	}
	tx.Commit()
	return t.SucceedWithResult(newResource)
}

func (p *RestHandle) HandlePut(store rest.ResourceStore, t *quark.Task) *quark.TaskResult {
	p.logger.Error("not support put cmd")
	c, _ := t.Cmds[0].(*rest.PutCmd)
	newRes, _ := c.NewResource.(*PersonalData)
	return t.Failed(errors.New("not support cmds"))
}

func (p *RestHandle) HandleDelete(store rest.ResourceStore, t *quark.Task) *quark.TaskResult {
	var err error
	if len(t.Cmds) > 1 {
		p.logger.Error("not support batch cmds")
		return t.Failed(errors.New("not support batch cmds"))
	}
	tx, _ := store.Begin()
	persional := []resource.PersonalData{}
	c, _ := t.Cmds[0].(*rest.DeleteCmd)
	err = tx.Fill(map[string]interface{}{"id": c.Id}, &persional)
	if err != nil {
		tx.RollBack()
		p.logger.Error("get persional data failed:" + err.Error())
		return t.Failed(errors.New("get persional data failed"))
	}
	if len(persional) != 1 {
		tx.RollBack()
		p.logger.Error("more than one record.")
		return t.Failed(errors.New("more than one record."))
	}

	_, err = tx.Delete(rest.ResourceType(c.ResourceType), map[string]interface{}{"id": c.Id})

	if err != nil {
		tx.RollBack()
		p.logger.Error("delete persional data failed:" + err.Error())
		return t.Failed(errors.New("delete persional data failed."))
	}
	tx.Commit()
	return t.SucceedWithResult(&persional[0])
}

func (p *RestHandle) SupportedResources() []rest.Resource {
	return []rest.Resource{
		&usermanager.Zdnsuser{},
		&resource.PersonalData{},
	}
}
