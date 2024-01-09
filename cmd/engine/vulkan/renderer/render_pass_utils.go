package renderer

import (
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	sealPass "github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/vulkan-go/vulkan"
)

func BeginVulkanRenderPass(renderer *VulkanRenderer, pass *sealPass.VulkanRenderPass, buffer *sealBuffer.VulkanCommandBuffer, imageIndex uint32) error {
	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  pass.Handle,
		Framebuffer: renderer.Swapchain.Framebuffers[imageIndex].Handle,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: renderer.Layout.Device.Physical.Capabilities.Surface.ImageExtent,
		},
		ClearValueCount: uint32(len(pass.Options.ClearValues)),
		PClearValues:    pass.Options.ClearValues,
	}
	vulkan.CmdBeginRenderPass(buffer.Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)

	return nil
}
