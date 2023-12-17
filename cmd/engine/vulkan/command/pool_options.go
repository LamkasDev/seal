package command

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandPoolOptions struct {
	CreateInfo vulkan.CommandPoolCreateInfo
}

func NewVulkanCommandPoolOptions(family uint32) (VulkanCommandPoolOptions, error) {
	options := VulkanCommandPoolOptions{
		CreateInfo: vulkan.CommandPoolCreateInfo{
			SType:            vulkan.StructureTypeCommandPoolCreateInfo,
			Flags:            vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateResetCommandBufferBit),
			QueueFamilyIndex: family,
		},
	}

	return options, nil
}
