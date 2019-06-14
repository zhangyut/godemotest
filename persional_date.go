package resource

type PersonalData struct {
	Id       string   `json:"id"`
	Zdnsuser string   `sql:"ownby" json:"zdnsuser"`
	Username string   `json:"username"`
	Addrs    []string `json:"addrs"`
}

func (p *PersonalData) Validate() error {
	return nil
}
