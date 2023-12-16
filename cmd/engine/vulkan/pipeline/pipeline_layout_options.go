package pipeline

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayoutOptions struct {
	CreateInfo vulkan.PipelineLayoutCreateInfo
}

func NewVulkanPipelineLayoutOptions() (VulkanPipelineLayoutOptions, error) {
	options := VulkanPipelineLayoutOptions{
		CreateInfo: vulkan.PipelineLayoutCreateInfo{
			SType: vulkan.StructureTypePipelineLayoutCreateInfo,
		},
	}

	return options, nil
}
