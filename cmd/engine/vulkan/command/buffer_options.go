package command

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandBufferOptions struct {
	AllocateInfo vulkan.CommandBufferAllocateInfo
}

func NewVulkanCommandBufferOptions(pool *VulkanCommandPool) VulkanCommandBufferOptions {
	options := VulkanCommandBufferOptions{
		AllocateInfo: vulkan.CommandBufferAllocateInfo{
			SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
			CommandPool:        pool.Handle,
			Level:              vulkan.CommandBufferLevelPrimary,
			CommandBufferCount: 1,
		},
	}

	return options
}
