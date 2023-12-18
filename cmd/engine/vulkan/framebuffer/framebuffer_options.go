package framebuffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/vulkan-go/vulkan"
)

type VulkanFramebufferOptions struct {
	CreateInfo vulkan.FramebufferCreateInfo
}

func NewVulkanFramebufferOptions(pass *pass.VulkanRenderPass, imageview *image.VulkanImageView, extent vulkan.Extent2D) VulkanFramebufferOptions {
	options := VulkanFramebufferOptions{
		CreateInfo: vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      pass.Handle,
			AttachmentCount: 1,
			PAttachments:    []vulkan.ImageView{imageview.Handle},
			Width:           extent.Width,
			Height:          extent.Height,
			Layers:          1,
		},
	}

	return options
}
