package renderer

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

func QueueSubmitVulkanRenderer(renderer *VulkanRenderer) error {
	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{renderer.RendererCommander.CurrentCommandBuffer.Handle},
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{renderer.RendererSyncer.ImageAvailableSemaphores[renderer.CurrentFrame].Handle},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{renderer.RendererSyncer.RenderFinishedSemaphores[renderer.CurrentFrame].Handle},
	}
	if res := vulkan.QueueSubmit(renderer.Layout.Device.Queues[uint32(renderer.Layout.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, renderer.RendererSyncer.InFlightFences[renderer.CurrentFrame].Handle); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func QueuePresentRenderer(renderer *VulkanRenderer, imageIndex uint32) error {
	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{renderer.RendererSyncer.RenderFinishedSemaphores[renderer.CurrentFrame].Handle},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{renderer.Swapchain.Handle},
		PImageIndices:      []uint32{imageIndex},
	}
	if res := vulkan.QueuePresent(renderer.Layout.Device.Queues[uint32(renderer.Layout.Device.Physical.Capabilities.Queue.PresentationIndex)], &presentInfo); res != vulkan.Success {
		switch res {
		case vulkan.ErrorOutOfDate:
		case vulkan.Suboptimal:
			if err := ResizeVulkanRenderer(renderer); err != nil {
				return err
			}
			return nil
		default:
			return vulkan.Error(res)
		}
	}

	return nil
}
