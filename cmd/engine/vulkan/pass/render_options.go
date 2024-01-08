package pass

import "github.com/vulkan-go/vulkan"

type VulkanRenderPassOptions struct {
	ColorAttachmentDescription vulkan.AttachmentDescription
	ColorAttachmentReference   vulkan.AttachmentReference
	DepthAttachmentDescription vulkan.AttachmentDescription
	DepthAttachmentReference   vulkan.AttachmentReference
	SubpassDescription         vulkan.SubpassDescription
	SubpassDependency          vulkan.SubpassDependency
	CreateInfo                 vulkan.RenderPassCreateInfo
}

func NewVulkanRenderPassOptions(format vulkan.Format) VulkanRenderPassOptions {
	options := VulkanRenderPassOptions{
		ColorAttachmentDescription: vulkan.AttachmentDescription{
			Format:         format,
			Samples:        vulkan.SampleCount1Bit,
			LoadOp:         vulkan.AttachmentLoadOpClear,
			StoreOp:        vulkan.AttachmentStoreOpStore,
			StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
			StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
			InitialLayout:  vulkan.ImageLayoutUndefined,
			FinalLayout:    vulkan.ImageLayoutPresentSrc,
		},
		ColorAttachmentReference: vulkan.AttachmentReference{
			Attachment: 0,
			Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
		},
		DepthAttachmentDescription: vulkan.AttachmentDescription{
			Format:         vulkan.FormatD32Sfloat,
			Samples:        vulkan.SampleCount1Bit,
			LoadOp:         vulkan.AttachmentLoadOpClear,
			StoreOp:        vulkan.AttachmentStoreOpDontCare,
			StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
			StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
			InitialLayout:  vulkan.ImageLayoutUndefined,
			FinalLayout:    vulkan.ImageLayoutDepthStencilAttachmentOptimal,
		},
		DepthAttachmentReference: vulkan.AttachmentReference{
			Attachment: 1,
			Layout:     vulkan.ImageLayoutDepthStencilAttachmentOptimal,
		},
		SubpassDependency: vulkan.SubpassDependency{
			SrcSubpass:    vulkan.SubpassExternal,
			DstSubpass:    0,
			SrcStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit | vulkan.PipelineStageEarlyFragmentTestsBit),
			DstStageMask:  vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit | vulkan.PipelineStageEarlyFragmentTestsBit),
			SrcAccessMask: 0,
			DstAccessMask: vulkan.AccessFlags(vulkan.AccessColorAttachmentWriteBit | vulkan.AccessDepthStencilAttachmentWriteBit),
		},
	}
	options.SubpassDescription = vulkan.SubpassDescription{
		PipelineBindPoint:       vulkan.PipelineBindPointGraphics,
		ColorAttachmentCount:    1,
		PColorAttachments:       []vulkan.AttachmentReference{options.ColorAttachmentReference},
		PDepthStencilAttachment: &options.DepthAttachmentReference,
	}
	options.CreateInfo = vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 2,
		PAttachments:    []vulkan.AttachmentDescription{options.ColorAttachmentDescription, options.DepthAttachmentDescription},
		SubpassCount:    1,
		PSubpasses:      []vulkan.SubpassDescription{options.SubpassDescription},
		DependencyCount: 1,
		PDependencies:   []vulkan.SubpassDependency{options.SubpassDependency},
	}

	return options
}
