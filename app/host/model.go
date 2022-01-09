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

//主机信息
type Host struct {
	Resource_hash string `json:"resource_hash"`
	Describe_hash string `json:"describe_hash"`
	*Resource
	*Describe
}

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

//定义厂商类型
type Vendor int

//使用iota枚举来定义厂商
const (
	ALI_CLOUD Vendor = iota
	TX_CLOUD
	HW_CLOUD
)

//主机统用信息
type Resource struct {
	Id     string `json:"id"`                         // 全局唯一Id
	Vendor Vendor `json:"vendor"`                     // 厂商
	Region string `json:"region" validate:"required"` // 地域
	Zone   string `json:"zone"`                       // 区域
	//使用13位的时间戳
	//为什么不用Datetime，如果使用数据时间,数据库会默认加上时区
	CreateAt    int64             `json:"create_at" validate:"required"`  // 创建时间
	ExpireAt    int64             `json:"expire_at"`                      // 过期时间
	Category    string            `json:"category"`                       // 种类
	Type        string            `json:"type"`                           // 规格
	InstanceId  string            `json:"instance_id"`                    // 实例ID
	Name        string            `json:"name" validate:"required"`       // 名称
	Description string            `json:"description"`                    // 描述
	Status      string            `json:"status" validate:"required"`     // 服务商中的状态
	Tags        map[string]string `json:"tags"`                           // 标签
	UpdateAt    int64             `json:"update_at"`                      // 更新时间
	SyncAt      int64             `json:"sync_at"`                        // 同步时间
	SyncAccount string            `json:"sync_accout"`                    // 同步的账号
	PublicIP    string            `json:"public_ip"`                      // 公网IP
	PrivateIP   string            `json:"private_ip" validate:"required"` // 内网IP
	PayType     string            `json:"pay_type"`                       // 实例付费方式
}

//主机详细信息
type Describe struct {
	CPU                     int    `json:"cpu" validate:"required"`    // 核数
	Memory                  int    `json:"memory" validate:"required"` // 内存
	GPUAmount               int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec                 string `json:"gpu_spec"`                   // GPU类型
	OSType                  string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName                  string `json:"os_name"`                    // 操作系统名称
	SerialNumber            string `json:"serial_number"`              // 序列号
	ImageID                 string `json:"image_id"`                   // 镜像ID
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"` // 公网出带宽最大值，单位为 Mbps
	InternetMaxBandwidthIn  int    `json:"internet_max_bandwidth_in"`  // 公网入带宽最大值，单位为 Mbps
	KeyPairName             string `json:"key_pair_name"`              // 秘钥对名称
	SecurityGroups          string `json:"security_groups"`            // 安全组  采用逗号分隔
}

//查询返回的集合定义
type HostSet struct {
	//数据总数
	Total int64 `json:"total"`
	//查询出的主机
	Items []*Host `json:"items"`
}

func NewDefaultHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}
