package descriptor

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorPoolOptions struct {
	CreateInfo vulkan.DescriptorPoolCreateInfo
}

func NewVulkanDescriptorPoolOptions() VulkanDescriptorPoolOptions {
	options := VulkanDescriptorPoolOptions{
		CreateInfo: vulkan.DescriptorPoolCreateInfo{
			SType:         vulkan.StructureTypeDescriptorPoolCreateInfo,
			PoolSizeCount: 1,
			PPoolSizes:    []vulkan.DescriptorPoolSize{{Type: vulkan.DescriptorTypeUniformBuffer, DescriptorCount: commonPipeline.MaxFramesInFlight}},
			MaxSets:       commonPipeline.MaxFramesInFlight,
		},
	}

	return options
}
