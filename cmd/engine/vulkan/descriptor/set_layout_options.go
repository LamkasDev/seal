package descriptor

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorSetLayoutOptions struct {
	CreateInfo vulkan.DescriptorSetLayoutCreateInfo
}

func NewVulkanDescriptorSetLayoutOptions() VulkanDescriptorSetLayoutOptions {
	options := VulkanDescriptorSetLayoutOptions{
		CreateInfo: vulkan.DescriptorSetLayoutCreateInfo{
			SType: vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
			PBindings: []vulkan.DescriptorSetLayoutBinding{
				{
					Binding:         0,
					DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
					DescriptorCount: 1,
					StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit),
				},
				{
					Binding:         1,
					DescriptorType:  vulkan.DescriptorTypeCombinedImageSampler,
					DescriptorCount: 1,
					StageFlags:      vulkan.ShaderStageFlags(vulkan.ShaderStageFragmentBit),
				},
			},
		},
	}
	options.CreateInfo.BindingCount = uint32(len(options.CreateInfo.PBindings))

	return options
}
