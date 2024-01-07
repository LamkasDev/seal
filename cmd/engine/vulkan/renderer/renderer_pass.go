package renderer

import (
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/vulkan-go/vulkan"
)

func BeginVulkanRenderPass(renderer *VulkanRenderer, buffer *sealBuffer.VulkanCommandBuffer, imageIndex uint32) error {
	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  renderer.RenderPass.Handle,
		Framebuffer: renderer.Swapchain.Framebuffers[imageIndex].Handle,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: renderer.Layout.Device.Physical.Capabilities.Surface.ImageExtent,
		},
		ClearValueCount: 1,
		PClearValues:    []vulkan.ClearValue{{1, 0, 0, 1}},
	}
	vulkan.CmdBeginRenderPass(buffer.Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)

	return nil
}
