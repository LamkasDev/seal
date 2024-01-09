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

	Shaders     []string
	ClearValues []vulkan.ClearValue
}
