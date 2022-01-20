package host

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

/*
	model用来定义主机数据结构和返回集合的数据结构
*/

var (
	validate = validator.New()
)



func NewDefaultHost() *Host {
	return &Host{
		Resource: &Resource{
			CreateAt: time.Now().UnixNano() / 1000000,
		},
		Describe: &Describe{},
	}
}

func (h *Host) Validate() error {
	return validate.Struct(h)
}

func (h *Host) Patch(res *Resource, desc *Describe) error {
	if res != nil {
		//将相同类型的结构体进行合并
		err := mergo.MergeWithOverwrite(h.Resource, res)
		if err != nil {
			return err
		}
	}

	if desc != nil {
		//将相同类型的结构体进行合并
		err := mergo.MergeWithOverwrite(h.Describe, desc)
		if err != nil {
			return err
		}
	}
	return nil
}

//go 1.17 允许获取毫秒
func (h *Host) Update(res *Resource, desc *Describe) {
	h.Resource = res
	h.Describe = desc
}


func NewDefaultHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}
