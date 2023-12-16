package pass

import "github.com/vulkan-go/vulkan"

type VulkanRenderPassOptions struct {
	ColorAttachmentDescription vulkan.AttachmentDescription
	ColorAttachmentReference   vulkan.AttachmentReference
	SubpassDescription         vulkan.SubpassDescription
	CreateInfo                 vulkan.RenderPassCreateInfo
}

func NewVulkanRenderPassOptions(format vulkan.Format) (VulkanRenderPassOptions, error) {
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
	}
	options.ColorAttachmentReference = vulkan.AttachmentReference{
		Attachment: 0,
		Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
	}
	options.SubpassDescription = vulkan.SubpassDescription{
		PipelineBindPoint:    vulkan.PipelineBindPointGraphics,
		ColorAttachmentCount: 1,
		PColorAttachments:    []vulkan.AttachmentReference{options.ColorAttachmentReference},
	}
	options.CreateInfo = vulkan.RenderPassCreateInfo{
		SType:           vulkan.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 1,
		PAttachments:    []vulkan.AttachmentDescription{options.ColorAttachmentDescription},
		SubpassCount:    1,
		PSubpasses:      []vulkan.SubpassDescription{options.SubpassDescription},
	}

	return options, nil
}
