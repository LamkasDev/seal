package descriptor

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/samber/lo"
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorSetOptions struct {
	AllocateInfo vulkan.DescriptorSetAllocateInfo
}

func NewVulkanDescriptorSetOptions(pool *VulkanDescriptorPool, setLayout *VulkanDescriptorSetLayout) VulkanDescriptorSetOptions {
	options := VulkanDescriptorSetOptions{
		AllocateInfo: vulkan.DescriptorSetAllocateInfo{
			SType:              vulkan.StructureTypeDescriptorSetAllocateInfo,
			DescriptorPool:     pool.Handle,
			DescriptorSetCount: 1,
			PSetLayouts: lo.RepeatBy(commonPipeline.MaxFramesInFlight, func(index int) vulkan.DescriptorSetLayout {
				return setLayout.Handle
			}),
		},
	}

	return options
}
