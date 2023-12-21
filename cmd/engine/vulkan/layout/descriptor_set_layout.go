package layout

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorSetLayoutOptions struct {
	CreateInfo vulkan.DescriptorSetLayoutCreateInfo
}

func NewVulkanDescriptorSetLayoutOptions() VulkanDescriptorSetLayoutOptions {
	options := VulkanDescriptorSetLayoutOptions{
		CreateInfo: vulkan.DescriptorSetLayoutCreateInfo{
			SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
			BindingCount: 1,
			PBindings: []vulkan.DescriptorSetLayoutBinding{
				{
					Binding:         0,
					DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
					DescriptorCount: 1,
					StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit),
				},
			},
		},
	}

	return options
}
