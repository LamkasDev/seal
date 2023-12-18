package command

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandPoolOptions struct {
	CreateInfo vulkan.CommandPoolCreateInfo
}

func NewVulkanCommandPoolOptions(queueFamilyIndex uint32) VulkanCommandPoolOptions {
	options := VulkanCommandPoolOptions{
		CreateInfo: vulkan.CommandPoolCreateInfo{
			SType:            vulkan.StructureTypeCommandPoolCreateInfo,
			Flags:            vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateResetCommandBufferBit),
			QueueFamilyIndex: queueFamilyIndex,
		},
	}

	return options
}
