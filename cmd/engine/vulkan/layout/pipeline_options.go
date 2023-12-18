package layout

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayoutOptions struct {
	CreateInfo vulkan.PipelineLayoutCreateInfo
}

func NewVulkanPipelineLayoutOptions() VulkanPipelineLayoutOptions {
	options := VulkanPipelineLayoutOptions{
		CreateInfo: vulkan.PipelineLayoutCreateInfo{
			SType: vulkan.StructureTypePipelineLayoutCreateInfo,
		},
	}

	return options
}
