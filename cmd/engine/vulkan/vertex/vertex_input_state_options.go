package vertex

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanVertexInputStateOptions struct {
	CreateInfo vulkan.PipelineVertexInputStateCreateInfo
}

func NewVertexInputStateOptions() VulkanVertexInputStateOptions {
	bindings := []vulkan.VertexInputBindingDescription{
		{Binding: 0, Stride: uint32(VulkanVertexSize), InputRate: vulkan.VertexInputRateVertex},
	}
	attributes := []vulkan.VertexInputAttributeDescription{
		{Binding: 0, Location: 0, Format: vulkan.FormatR32g32b32Sfloat, Offset: uint32(VulkanVertexPositionOffset)},
		{Binding: 0, Location: 1, Format: vulkan.FormatR32g32b32Sfloat, Offset: uint32(VulkanVertexColorOffset)},
		{Binding: 0, Location: 2, Format: vulkan.FormatR32g32Sfloat, Offset: uint32(VulkanVertexTexCoordOffset)},
	}
	options := VulkanVertexInputStateOptions{
		CreateInfo: vulkan.PipelineVertexInputStateCreateInfo{
			SType:                           vulkan.StructureTypePipelineVertexInputStateCreateInfo,
			VertexBindingDescriptionCount:   uint32(len(bindings)),
			PVertexBindingDescriptions:      bindings,
			VertexAttributeDescriptionCount: uint32(len(attributes)),
			PVertexAttributeDescriptions:    attributes,
		},
	}

	return options
}
