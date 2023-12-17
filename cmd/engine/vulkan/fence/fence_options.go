package fence

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanFenceOptions struct {
	CreateInfo vulkan.FenceCreateInfo
}

func NewVulkanFenceOptions(flags vulkan.FenceCreateFlags) (VulkanFenceOptions, error) {
	options := VulkanFenceOptions{
		CreateInfo: vulkan.FenceCreateInfo{
			SType: vulkan.StructureTypeFenceCreateInfo,
			Flags: flags,
		},
	}

	return options, nil
}
