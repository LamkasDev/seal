package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline_layout"
	sealShader "github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineOptions struct {
	DynamicState              vulkan.PipelineDynamicStateCreateInfo
	VertexInputStateOptions   vertex.VulkanVertexInputStateOptions
	InputAssemblyState        vulkan.PipelineInputAssemblyStateCreateInfo
	ViewportState             vulkan.PipelineViewportStateCreateInfo
	RasterizationState        vulkan.PipelineRasterizationStateCreateInfo
	MultisampleState          vulkan.PipelineMultisampleStateCreateInfo
	ColorBlendAttachmentState vulkan.PipelineColorBlendAttachmentState
	ColorBlendState           vulkan.PipelineColorBlendStateCreateInfo
	CreateInfo                vulkan.GraphicsPipelineCreateInfo
}

func NewVulkanPipelineOptions(layout *pipeline_layout.VulkanPipelineLayout, viewport *viewport.VulkanViewport, pass *pass.VulkanRenderPass, shader *sealShader.VulkanShader) VulkanPipelineOptions {
	options := VulkanPipelineOptions{
		DynamicState: vulkan.PipelineDynamicStateCreateInfo{
			SType:             vulkan.StructureTypePipelineDynamicStateCreateInfo,
			DynamicStateCount: 2,
			PDynamicStates:    []vulkan.DynamicState{vulkan.DynamicStateViewport, vulkan.DynamicStateScissor},
		},
		VertexInputStateOptions: vertex.NewVertexInputStateOptions(),
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
			FrontFace:               vulkan.FrontFaceCounterClockwise,
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

	options.CreateInfo = vulkan.GraphicsPipelineCreateInfo{
		SType:      vulkan.StructureTypeGraphicsPipelineCreateInfo,
		StageCount: 2,
		PStages: []vulkan.PipelineShaderStageCreateInfo{
			shader.Vertex.Stage, shader.Fragment.Stage,
		},
		PVertexInputState:   &options.VertexInputStateOptions.CreateInfo,
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
