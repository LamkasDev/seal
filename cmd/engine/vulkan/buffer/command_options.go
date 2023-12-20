package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandBufferOptions struct {
	AllocateInfo vulkan.CommandBufferAllocateInfo
}

func NewVulkanCommandBufferOptions(pool *command.VulkanCommandPool) VulkanCommandBufferOptions {
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
