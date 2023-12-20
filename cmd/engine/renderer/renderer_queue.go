package renderer

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

func QueueSubmitVulkanRenderer(renderer *Renderer) error {
	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{renderer.Pipeline.Syncer.ImageAvailableSemaphores[renderer.Pipeline.CurrentFrame].Handle},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{renderer.Pipeline.Commander.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{renderer.Pipeline.Syncer.RenderFinishedSemaphores[renderer.Pipeline.CurrentFrame].Handle},
	}
	if res := vulkan.QueueSubmit(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, renderer.Pipeline.Syncer.InFlightFences[renderer.Pipeline.CurrentFrame].Handle); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func QueuePresentRenderer(renderer *Renderer, imageIndex uint32) error {
	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{renderer.Pipeline.Syncer.RenderFinishedSemaphores[renderer.Pipeline.CurrentFrame].Handle},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{renderer.Swapchain.Handle},
		PImageIndices:      []uint32{imageIndex},
	}
	if res := vulkan.QueuePresent(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.PresentationIndex)], &presentInfo); res != vulkan.Success {
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
