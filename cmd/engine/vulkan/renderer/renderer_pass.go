package renderer

import (
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/vulkan-go/vulkan"
)

func BeginVulkanRenderPass(renderer *VulkanRenderer, buffer *sealBuffer.VulkanCommandBuffer, imageIndex uint32) error {
	colors := []vulkan.ClearValue{{}, {}}
	colors[0].SetColor([]float32{0, 0, 0, 1})
	colors[1].SetDepthStencil(1, 0)

	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  renderer.RenderPass.Handle,
		Framebuffer: renderer.Swapchain.Framebuffers[imageIndex].Handle,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: renderer.Layout.Device.Physical.Capabilities.Surface.ImageExtent,
		},
		ClearValueCount: uint32(len(colors)),
		PClearValues:    colors,
	}
	vulkan.CmdBeginRenderPass(buffer.Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)

	return nil
}
