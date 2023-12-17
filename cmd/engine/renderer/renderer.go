package renderer

import (
	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
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
	renderer.Window, err = window.NewWindow(window.NewWindowOptions("Test", constants.DefaultResolution))
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
	if err := sealVulkan.InitializeVulkanInstanceDevices(&renderer.VulkanInstance, renderer.Window.Handle, &renderer.Surface); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.ShaderContainer, err = shader.NewVulkanShaderContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Pipeline, err = pipeline.NewVulkanPipeline(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.ShaderContainer); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.Pipeline, &renderer.Surface); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func RunRenderer(renderer *Renderer) error {
	if res := vulkan.WaitForFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.InFlightFence.Handle}, vulkan.Bool32(1), vulkan.MaxUint64); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	if res := vulkan.ResetFences(renderer.Pipeline.Device.Handle, 1, []vulkan.Fence{renderer.Pipeline.InFlightFence.Handle}); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	var imageIndex uint32
	if res := vulkan.AcquireNextImage(renderer.Pipeline.Device.Handle, renderer.Swapchain.Handle, vulkan.MaxUint64, renderer.Pipeline.SemaphoreImageAvailable.Handle, nil, &imageIndex); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	if res := vulkan.ResetCommandBuffer(renderer.Pipeline.CommandBuffer.Handle, 0); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	commandBufferBeginInfo := vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: vulkan.CommandBufferUsageFlags(vulkan.CommandBufferUsageOneTimeSubmitBit),
	}
	if res := vulkan.BeginCommandBuffer(renderer.Pipeline.CommandBuffer.Handle, &commandBufferBeginInfo); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
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
	vulkan.CmdBeginRenderPass(renderer.Pipeline.CommandBuffer.Handle, &renderPassBeginInfo, vulkan.SubpassContentsInline)
	vulkan.CmdBindPipeline(renderer.Pipeline.CommandBuffer.Handle, vulkan.PipelineBindPointGraphics, renderer.Pipeline.Handle)
	vulkan.CmdSetViewport(renderer.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Viewport{
		{
			X:        0,
			Y:        0,
			Width:    float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Width),
			Height:   float32(renderer.Swapchain.Options.CreateInfo.ImageExtent.Height),
			MinDepth: 0,
			MaxDepth: 0,
		},
	})
	vulkan.CmdSetScissor(renderer.Pipeline.CommandBuffer.Handle, 0, 1, []vulkan.Rect2D{{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: renderer.Swapchain.Options.CreateInfo.ImageExtent}})
	vulkan.CmdDraw(renderer.Pipeline.CommandBuffer.Handle, 3, 1, 0, 0)
	vulkan.CmdEndRenderPass(renderer.Pipeline.CommandBuffer.Handle)
	if res := vulkan.EndCommandBuffer(renderer.Pipeline.CommandBuffer.Handle); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	submitInfo := vulkan.SubmitInfo{
		SType:                vulkan.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vulkan.Semaphore{renderer.Pipeline.SemaphoreImageAvailable.Handle},
		PWaitDstStageMask:    []vulkan.PipelineStageFlags{vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vulkan.CommandBuffer{renderer.Pipeline.CommandBuffer.Handle},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vulkan.Semaphore{renderer.Pipeline.SemaphoreRenderFinished.Handle},
	}
	if res := vulkan.QueueSubmit(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, renderer.Pipeline.InFlightFence.Handle); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	presentInfo := vulkan.PresentInfo{
		SType:              vulkan.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vulkan.Semaphore{renderer.Pipeline.SemaphoreRenderFinished.Handle},
		SwapchainCount:     1,
		PSwapchains:        []vulkan.Swapchain{renderer.Swapchain.Handle},
		PImageIndices:      []uint32{imageIndex},
	}
	if res := vulkan.QueuePresent(renderer.Pipeline.Device.Queues[uint32(renderer.Pipeline.Device.Physical.Capabilities.Queue.PresentationIndex)], &presentInfo); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	return nil
}

func FreeRenderer(renderer *Renderer) error {
	vulkan.DeviceWaitIdle(renderer.Pipeline.Device.Handle)
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
