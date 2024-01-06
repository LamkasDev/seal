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
			SType: vulkan.StructureTypeDescriptorPoolCreateInfo,
			PPoolSizes: []vulkan.DescriptorPoolSize{
				{
					Type:            vulkan.DescriptorTypeUniformBuffer,
					DescriptorCount: commonPipeline.MaxFramesInFlight,
				},
				{
					Type:            vulkan.DescriptorTypeCombinedImageSampler,
					DescriptorCount: commonPipeline.MaxFramesInFlight,
				},
			},
			MaxSets: commonPipeline.MaxFramesInFlight,
		},
	}
	options.CreateInfo.PoolSizeCount = uint32(len(options.CreateInfo.PPoolSizes))

	return options
}
