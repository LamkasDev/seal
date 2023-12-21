package layout

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayoutOptions struct {
	DescriptorSet *VulkanDescriptorSetLayout
	CreateInfo    vulkan.PipelineLayoutCreateInfo
}

func NewVulkanPipelineLayoutOptions(descriptorSet *VulkanDescriptorSetLayout) VulkanPipelineLayoutOptions {
	options := VulkanPipelineLayoutOptions{
		CreateInfo: vulkan.PipelineLayoutCreateInfo{
			SType:          vulkan.StructureTypePipelineLayoutCreateInfo,
			SetLayoutCount: 1,
			PSetLayouts:    []vulkan.DescriptorSetLayout{descriptorSet.Handle},
		},
	}

	return options
}
