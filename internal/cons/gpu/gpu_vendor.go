package gpu

import (
	"strings"
)

// GpuVendor 表示显卡供应商类型
type GpuVendor string

const (
	// NVIDIA 英伟达显卡
	NVIDIA GpuVendor = "NVIDIA (0x10de)"

	// AMD AMD显卡
	AMD GpuVendor = "AMD"

	// INTEL 英特尔显卡
	INTEL GpuVendor = "Intel Corporation (0x8086)"

	// UNKNOWN 未知显卡供应商
	UNKNOWN GpuVendor = "UNKNOWN"
)

// String 返回显卡供应商的字符串表示
func (v GpuVendor) String() string {
	return string(v)
}

// GetVendor 返回显卡供应商的名称
func (v GpuVendor) GetVendor() string {
	return string(v)
}

// GetByVendor 根据显卡供应商名称获取对应的枚举值
func GetByVendor(vendor string) GpuVendor {
	if vendor == "" {
		return UNKNOWN
	}

	vendorUpper := strings.ToUpper(vendor)

	// 使用switch语句使代码更简洁
	for _, v := range []GpuVendor{NVIDIA, AMD, INTEL} {
		if strings.Contains(vendorUpper, strings.ToUpper(string(v))) {
			return v
		}
	}

	return UNKNOWN
}
