package renderer

import (
	"github.com/vulkan-go/vulkan"
)

func BeginVulkanRenderPass(renderer *Renderer, imageIndex uint32) error {
	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  renderer.Pipeline.RenderPass.Handle,
		Framebuffer: renderer.Swapchain.Framebuffers[imageIndex].Handle,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: renderer.Pipeline.Device.Physical.Capabilities.Surface.ImageExtent,
		},
		ClearValueCount: 1,
		PClearValues:    []vulkan.ClearValue{{0, 0, 0, 1}},
	}
	vulkan.CmdBeginRenderPass(renderer.Pipeline.Commander.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)

	return nil
}
