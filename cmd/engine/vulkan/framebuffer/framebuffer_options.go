package framebuffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/vulkan-go/vulkan"
)

type VulkanFramebufferOptions struct {
	CreateInfo vulkan.FramebufferCreateInfo
}

func NewVulkanFramebufferOptions(device *logical.VulkanLogicalDevice, pass *pass.VulkanRenderPass, imageview *image.VulkanImageView) VulkanFramebufferOptions {
	options := VulkanFramebufferOptions{
		CreateInfo: vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      pass.Handle,
			AttachmentCount: 1,
			PAttachments:    []vulkan.ImageView{imageview.Handle},
			Width:           device.Physical.Capabilities.Surface.ImageExtent.Width,
			Height:          device.Physical.Capabilities.Surface.ImageExtent.Height,
			Layers:          1,
		},
	}

	return options
}
