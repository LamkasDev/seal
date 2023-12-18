package renderer

import (
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/device"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/swapchain"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type Renderer struct {
	VulkanInstance  sealVulkan.VulkanInstance
	Window          window.Window
	Surface         vulkan.Surface
	ShaderContainer shader.VulkanShaderContainer
	Pipeline        pipeline.VulkanPipeline
	Swapchain       swapchain.VulkanSwapchain
}

func NewRenderer() (Renderer, error) {
	var err error
	renderer := Renderer{}

	progress.AdvanceLoading()
	renderer.VulkanInstance, err = sealVulkan.NewVulkanInstance()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	renderer.Window, err = window.NewWindow()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	var surfaceRaw uintptr
	if surfaceRaw, err = renderer.Window.Handle.CreateWindowSurface(renderer.VulkanInstance.Handle, nil); err != nil {
		return renderer, err
	}
	renderer.Surface = vulkan.Surface(vulkan.SurfaceFromPointer(surfaceRaw))

	progress.AdvanceLoading()
	if err := sealVulkan.InitializeVulkanInstanceDevices(&renderer.VulkanInstance, &renderer.Window, &renderer.Surface); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.ShaderContainer, err = shader.NewVulkanShaderContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Pipeline, err = pipeline.NewVulkanPipeline(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.ShaderContainer, &renderer.Window); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.Pipeline, &renderer.Surface, nil); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func UpdateVulkanRenderer(renderer *Renderer) error {
	var err error
	if err = swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err = device.UpdateVulkanInstanceDevices(&renderer.VulkanInstance.Devices); err != nil {
		return err
	}
	if err = pipeline.UpdateVulkanPipeline(&renderer.Pipeline); err != nil {
		return err
	}
	renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.Pipeline, &renderer.Surface, nil)
	if err != nil {
		return err
	}

	return nil
}

func AcquireNextImageRenderer(renderer *Renderer) (uint32, error) {
	var imageIndex uint32
	if res := vulkan.AcquireNextImage(renderer.Pipeline.Device.Handle, renderer.Swapchain.Handle, vulkan.MaxUint64, renderer.Pipeline.ImageAvailableSemaphores[renderer.Pipeline.CurrentFrame].Handle, nil, &imageIndex); res != vulkan.Success && res != vulkan.Suboptimal {
		switch res {
		case vulkan.ErrorOutOfDate:
			if err := UpdateVulkanRenderer(renderer); err != nil {
				return vulkan.MaxUint32, err
			}
			return AcquireNextImageRenderer(renderer)
		default:
			return vulkan.MaxUint32, vulkan.Error(res)
		}
	}

	return imageIndex, nil
}

func QueuePresentRenderer(renderer *Renderer, presentInfo vulkan.PresentInfo) error {
	if res := vulkan.QueuePresent(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.PresentationIndex)], &presentInfo); res != vulkan.Success {
		switch res {
		case vulkan.ErrorOutOfDate:
		case vulkan.Suboptimal:
			if err := UpdateVulkanRenderer(renderer); err != nil {
				return err
			}
			return nil
		default:
			return vulkan.Error(res)
		}
	}

	return nil
}

func RunRenderer(renderer *Renderer) error {
	var err error
	if res := vulkan.WaitForFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.InFlightFences[renderer.Pipeline.CurrentFrame].Handle}, vulkan.Bool32(1), vulkan.MaxUint64); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.ResetFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.InFlightFences[renderer.Pipeline.CurrentFrame].Handle}); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	var imageIndex uint32
	if imageIndex, err = AcquireNextImageRenderer(renderer); err != nil {
		logger.DefaultLogger.Error(err)
		return err
	}

	if res := vulkan.ResetCommandBuffer(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	commandBufferBeginInfo := vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: vulkan.CommandBufferUsageFlags(vulkan.CommandBufferUsageOneTimeSubmitBit),
	}
	if res := vulkan.BeginCommandBuffer(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, &commandBufferBeginInfo); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	renderPassBeginInfo := vulkan.RenderPassBeginInfo{
		SType:       vulkan.StructureTypeRenderPassBeginInfo,
		RenderPass:  renderer.Pipeline.RenderPass.Handle,
		Framebuffer: renderer.Swapchain.Framebuffers[imageIndex].Handle,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: renderer.Swapchain.Options.CreateInfo.ImageExtent,
		},
		ClearValueCount: 1,
		PClearValues:    []vulkan.ClearValue{{0, 0, 0, 1}},
	}
	vulkan.CmdBeginRenderPass(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)
	vulkan.CmdBindPipeline(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, vulkan.PipelineBindPointGraphics, renderer.Pipeline.Handle)
	vulkan.CmdSetViewport(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, 0, 1, []vulkan.Viewport{
		{
			X:        0,
			Y:        0,
			Width:    float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Width),
			Height:   float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Height),
			MinDepth: 0,
			MaxDepth: 0,
		},
	})
	vulkan.CmdSetScissor(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, 0, 1, []vulkan.Rect2D{{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: renderer.Swapchain.Options.CreateInfo.ImageExtent}})
	vulkan.CmdDraw(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle, 3, 1, 0, 0)
	vulkan.CmdEndRenderPass(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle)
	if res := vulkan.EndCommandBuffer(renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{renderer.Pipeline.ImageAvailableSemaphores[renderer.Pipeline.CurrentFrame].Handle},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{renderer.Pipeline.CommandBuffers[renderer.Pipeline.CurrentFrame].Handle},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{renderer.Pipeline.RenderFinishedSemaphores[renderer.Pipeline.CurrentFrame].Handle},
	}
	if res := vulkan.QueueSubmit(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, renderer.Pipeline.InFlightFences[renderer.Pipeline.CurrentFrame].Handle); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{renderer.Pipeline.RenderFinishedSemaphores[renderer.Pipeline.CurrentFrame].Handle},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{renderer.Swapchain.Handle},
		PImageIndices:      []uint32{imageIndex},
	}
	if err = QueuePresentRenderer(renderer, presentInfo); err != nil {
		logger.DefaultLogger.Error(err)
		return err
	}

	renderer.Pipeline.CurrentFrame = (renderer.Pipeline.CurrentFrame + 1) % pipeline.MaxFramesInFlight
	return nil
}

func FreeRenderer(renderer *Renderer) error {
	if err := swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	if err := pipeline.FreeVulkanPipeline(&renderer.Pipeline); err != nil {
		return err
	}
	if err := shader.FreeVulkanShaderContainer(&renderer.ShaderContainer); err != nil {
		return err
	}
	vulkan.DestroySurface(renderer.VulkanInstance.Handle, renderer.Surface, nil)
	if err := window.FreeWindow(&renderer.Window); err != nil {
		return err
	}

	return sealVulkan.FreeVulkanInstance(&renderer.VulkanInstance)
}
