package vertex

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanVertexInputStateOptions struct {
	CreateInfo vulkan.PipelineVertexInputStateCreateInfo
}

func NewVertexInputStateOptions() VulkanVertexInputStateOptions {
	options := VulkanVertexInputStateOptions{
		CreateInfo: vulkan.PipelineVertexInputStateCreateInfo{
			SType:                         vulkan.StructureTypePipelineVertexInputStateCreateInfo,
			VertexBindingDescriptionCount: 1,
			PVertexBindingDescriptions: []vulkan.VertexInputBindingDescription{
				{Binding: 0, Stride: uint32(VulkanVertexSize), InputRate: vulkan.VertexInputRateVertex},
			},
			VertexAttributeDescriptionCount: 1,
			PVertexAttributeDescriptions: []vulkan.VertexInputAttributeDescription{
				{Binding: 0, Location: 0, Format: vulkan.FormatR32g32Sfloat, Offset: uint32(VulkanVertexPositionOffset)},
				{Binding: 0, Location: 1, Format: vulkan.FormatR32g32b32Sfloat, Offset: uint32(VulkanVertexColorOffset)},
			},
		},
	}

	return options
}
