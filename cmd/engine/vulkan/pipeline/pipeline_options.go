package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/vulkan-go/vulkan"
	"golang.org/x/exp/maps"
)

type VulkanPipelineOptions struct {
	DynamicState              vulkan.PipelineDynamicStateCreateInfo
	VertexInputState          vulkan.PipelineVertexInputStateCreateInfo
	InputAssemblyState        vulkan.PipelineInputAssemblyStateCreateInfo
	ViewportState             vulkan.PipelineViewportStateCreateInfo
	RasterizationState        vulkan.PipelineRasterizationStateCreateInfo
	MultisampleState          vulkan.PipelineMultisampleStateCreateInfo
	ColorBlendAttachmentState vulkan.PipelineColorBlendAttachmentState
	ColorBlendState           vulkan.PipelineColorBlendStateCreateInfo
	CreateInfo                vulkan.GraphicsPipelineCreateInfo
}

func NewVulkanPipelineOptions(layout *layout.VulkanPipelineLayout, viewport *viewport.VulkanViewport, pass *pass.VulkanRenderPass, container *shader.VulkanShaderContainer) VulkanPipelineOptions {
	options := VulkanPipelineOptions{
		DynamicState: vulkan.PipelineDynamicStateCreateInfo{
			SType:             vulkan.StructureTypePipelineDynamicStateCreateInfo,
			DynamicStateCount: 2,
			PDynamicStates:    []vulkan.DynamicState{vulkan.DynamicStateViewport, vulkan.DynamicStateScissor},
		},
		VertexInputState: vulkan.PipelineVertexInputStateCreateInfo{
			SType:                           vulkan.StructureTypePipelineVertexInputStateCreateInfo,
			VertexBindingDescriptionCount:   0,
			VertexAttributeDescriptionCount: 0,
		},
		InputAssemblyState: vulkan.PipelineInputAssemblyStateCreateInfo{
			SType:                  vulkan.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology:               vulkan.PrimitiveTopologyTriangleList,
			PrimitiveRestartEnable: vulkan.False,
		},
		ViewportState: vulkan.PipelineViewportStateCreateInfo{
			SType:         vulkan.StructureTypePipelineViewportStateCreateInfo,
			ViewportCount: 1,
			PViewports:    []vulkan.Viewport{viewport.Viewport},
			ScissorCount:  1,
			PScissors:     []vulkan.Rect2D{viewport.Scissor},
		},
		RasterizationState: vulkan.PipelineRasterizationStateCreateInfo{
			SType:                   vulkan.StructureTypePipelineRasterizationStateCreateInfo,
			DepthClampEnable:        vulkan.Bool32(0),
			RasterizerDiscardEnable: vulkan.Bool32(0),
			PolygonMode:             vulkan.PolygonModeFill,
			LineWidth:               1,
			CullMode:                vulkan.CullModeFlags(vulkan.CullModeBackBit),
			FrontFace:               vulkan.FrontFaceClockwise,
			DepthBiasEnable:         vulkan.Bool32(0),
		},
		MultisampleState: vulkan.PipelineMultisampleStateCreateInfo{
			SType:                vulkan.StructureTypePipelineMultisampleStateCreateInfo,
			SampleShadingEnable:  vulkan.Bool32(0),
			RasterizationSamples: vulkan.SampleCount1Bit,
		},
		ColorBlendAttachmentState: vulkan.PipelineColorBlendAttachmentState{
			ColorWriteMask: vulkan.ColorComponentFlags(vulkan.ColorComponentRBit | vulkan.ColorComponentGBit | vulkan.ColorComponentBBit | vulkan.ColorComponentABit),
			BlendEnable:    vulkan.Bool32(0),
		},
	}
	options.ColorBlendState = vulkan.PipelineColorBlendStateCreateInfo{
		SType:           vulkan.StructureTypePipelineColorBlendStateCreateInfo,
		LogicOpEnable:   vulkan.Bool32(0),
		AttachmentCount: 1,
		PAttachments:    []vulkan.PipelineColorBlendAttachmentState{options.ColorBlendAttachmentState},
	}

	shaders := maps.Values(container.Shaders)
	stagesCount := len(shaders) * 2
	stages := make([]vulkan.PipelineShaderStageCreateInfo, stagesCount)
	for i := 0; i < stagesCount; i += 2 {
		stages[i] = shaders[i].Vertex.Stage
		stages[i+1] = shaders[i].Fragment.Stage
	}

	options.CreateInfo = vulkan.GraphicsPipelineCreateInfo{
		SType:               vulkan.StructureTypeGraphicsPipelineCreateInfo,
		StageCount:          uint32(stagesCount),
		PStages:             stages,
		PVertexInputState:   &options.VertexInputState,
		PInputAssemblyState: &options.InputAssemblyState,
		PDynamicState:       &options.DynamicState,
		PViewportState:      &options.ViewportState,
		PRasterizationState: &options.RasterizationState,
		PMultisampleState:   &options.MultisampleState,
		PColorBlendState:    &options.ColorBlendState,
		Layout:              layout.Handle,
		RenderPass:          pass.Handle,
		Subpass:             0,
	}

	return options
}
