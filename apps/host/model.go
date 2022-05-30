package host

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type HostSet struct {
	Total int     `json:"total"`
	Items []*Host `json:"items"`
}

type Host struct {
	// 资源公共属性
	*Resource
	// 资源独立属性
	*Describe
}

func NewHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}

func (h *Host) Validate() error {

	return validate.Struct(h)
}

// 填充默认值
func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

// Host对象全量更新
func (h *Host) Put(obj *Host) error {
	if obj.Id != h.Id { // 全局ID不允许更新
		return fmt.Errorf("Id not allowed to update, Permission Denied!")
	}
	*h.Resource = *obj.Resource
	*h.Describe = *obj.Describe

	return nil
}

// Host对象局部更新
func (h *Host) Patch(obj *Host) error {
	if obj.Name != "" {
		h.Name = obj.Name
	}

	if obj.CPU != 0 {
		h.CPU = obj.CPU
	}

	return nil
}

type Vendor int

const (
	// IDC, 枚举的默认值
	PRIVATE_IDC Vendor = iota
	// 阿里云
	ALIYUN
	// 腾讯云
	TXYUN
)

// 资源公共属性
type Resource struct {
	Id          string            `json:"id"  validate:"required"`     // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                      // 厂商
	Region      string            `json:"region"  validate:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                   // 创建时间
	ExpireAt    int64             `json:"expire_at"`                   // 过期时间
	Type        string            `json:"type"  validate:"required"`   // 规格
	Name        string            `json:"name"  validate:"required"`   // 名称
	Description string            `json:"description"`                 // 描述
	Status      string            `json:"status"`                      // 服务商中的状态
	Tags        map[string]string `json:"tags"`                        // 标签
	UpdateAt    int64             `json:"update_at"`                   // 更新时间
	SyncAt      int64             `json:"sync_at"`                     // 同步时间
	Account     string            `json:"accout"`                      // 资源的所属账号
	PublicIP    string            `json:"public_ip"`                   // 公网IP
	PrivateIP   string            `json:"private_ip"`                  // 内网IP
}

// 资源独立属性
type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

// 查询主机请求
type QueryHostRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	KeyWords   string `json:"kws"`
}

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func NewQueryHostFromHTTP(r *http.Request) *QueryHostRequest {
	req := NewQueryHostRequest()
	// query string
	qs := r.URL.Query()

	pss := qs.Get("page_size")
	if pss != "" {
		req.PageSize, _ = strconv.Atoi(pss)
	}

	pns := qs.Get("page_number")
	if pns != "" {
		req.PageNumber, _ = strconv.Atoi(pns)
	}

	req.KeyWords = qs.Get("kws")

	return req
}

func (req *QueryHostRequest) OffSet() int64 {
	return int64((req.PageNumber - 1) * req.PageSize)
}

func (req *QueryHostRequest) GetPageSize() uint {
	return uint(req.PageSize)
}

type DescribeHostRequest struct {
	Id string
}

func NewDescribeHostRequestWithId(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		Id: id,
	}
}

type UPDATE_MODE string

const (
	// 全量更新
	UPDATE_MODE_PUT UPDATE_MODE = "put"
	// 局部更新
	UPDATE_MODE_PATCH UPDATE_MODE = "patch"
)

// 更新主机信息请求
type UpdateHostRequest struct {
	UpdateMode UPDATE_MODE `json:"update_mode"`
	*Host
}

func NewPutUpdateRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PUT,
		Host:       h,
	}
}

func NewPatchUpdateRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PUT,
		Host:       h,
	}
}

// 删除主机请求
type DeleteHostRequest struct {
	Id string
}

func NewDeleteHostRequestWithId(id string) *DeleteHostRequest {
	return &DeleteHostRequest{
		Id: id,
	}
}

