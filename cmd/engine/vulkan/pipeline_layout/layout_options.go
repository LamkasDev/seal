package pipeline_layout

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayoutOptions struct {
	DescriptorSet *descriptor.VulkanDescriptorSetLayout
	CreateInfo    vulkan.PipelineLayoutCreateInfo
}

func NewVulkanPipelineLayoutOptions(descriptorSet *descriptor.VulkanDescriptorSetLayout) VulkanPipelineLayoutOptions {
	options := VulkanPipelineLayoutOptions{
		CreateInfo: vulkan.PipelineLayoutCreateInfo{
			SType:          vulkan.StructureTypePipelineLayoutCreateInfo,
			SetLayoutCount: 1,
			PSetLayouts:    []vulkan.DescriptorSetLayout{descriptorSet.Handle},
		},
	}

	return options
}
