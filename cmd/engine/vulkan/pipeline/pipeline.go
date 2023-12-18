package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/fence"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/semaphore"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

const MaxFramesInFlight = 2

type VulkanPipeline struct {
	Handle                   vulkan.Pipeline
	Device                   *logical.VulkanLogicalDevice
	Container                *shader.VulkanShaderContainer
	Window                   *window.Window
	Viewport                 viewport.VulkanViewport
	Layout                   layout.VulkanPipelineLayout
	RenderPass               pass.VulkanRenderPass
	CommandPool              command.VulkanCommandPool
	CommandBuffers           []command.VulkanCommandBuffer
	ImageAvailableSemaphores []semaphore.VulkanSemaphore
	RenderFinishedSemaphores []semaphore.VulkanSemaphore
	InFlightFences           []fence.VulkanFence
	CurrentFrame             uint32
	Options                  VulkanPipelineOptions
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, container *shader.VulkanShaderContainer, cwindow *window.Window) (VulkanPipeline, error) {
	var err error
	pipeline := VulkanPipeline{
		Device:                   device,
		Container:                container,
		Window:                   cwindow,
		Viewport:                 viewport.NewVulkanViewport(window.GetWindowImageExtent(cwindow)),
		CommandBuffers:           make([]command.VulkanCommandBuffer, MaxFramesInFlight),
		ImageAvailableSemaphores: make([]semaphore.VulkanSemaphore, MaxFramesInFlight),
		RenderFinishedSemaphores: make([]semaphore.VulkanSemaphore, MaxFramesInFlight),
		InFlightFences:           make([]fence.VulkanFence, MaxFramesInFlight),
	}

	if pipeline.Layout, err = layout.NewVulkanPipelineLayout(device); err != nil {
		return pipeline, err
	}
	if pipeline.RenderPass, err = pass.NewVulkanRenderPass(device, device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format); err != nil {
		return pipeline, err
	}
	if pipeline.CommandPool, err = command.NewVulkanCommandPool(device, uint32(device.Physical.Capabilities.Queue.GraphicsIndex)); err != nil {
		return pipeline, err
	}
	for i := 0; i < MaxFramesInFlight; i++ {
		if pipeline.CommandBuffers[i], err = command.NewVulkanCommandBuffer(device, &pipeline.CommandPool); err != nil {
			return pipeline, err
		}
		if pipeline.ImageAvailableSemaphores[i], err = semaphore.NewVulkanSemaphore(device); err != nil {
			return pipeline, err
		}
		if pipeline.RenderFinishedSemaphores[i], err = semaphore.NewVulkanSemaphore(device); err != nil {
			return pipeline, err
		}
		if pipeline.InFlightFences[i], err = fence.NewVulkanFence(device, vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit)); err != nil {
			return pipeline, err
		}
	}
	pipeline.Options = NewVulkanPipelineOptions(&pipeline.Layout, &pipeline.Viewport, &pipeline.RenderPass, container)

	vulkanPipelines := make([]vulkan.Pipeline, 1)
	if res := vulkan.CreateGraphicsPipelines(device.Handle, nil, 1, []vulkan.GraphicsPipelineCreateInfo{pipeline.Options.CreateInfo}, nil, vulkanPipelines); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pipeline, vulkan.Error(res)
	}
	pipeline.Handle = vulkanPipelines[0]
	logger.DefaultLogger.Debug("created new vulkan pipeline")

	return pipeline, nil
}

func UpdateVulkanPipeline(pipeline *VulkanPipeline) error {
	pipeline.Viewport = viewport.NewVulkanViewport(window.GetWindowImageExtent(pipeline.Window))
	return nil
}

func FreeVulkanPipeline(pipeline *VulkanPipeline) error {
	for i := 0; i < MaxFramesInFlight; i++ {
		if err := fence.FreeVulkanFence(&pipeline.InFlightFences[i]); err != nil {
			return err
		}
		if err := semaphore.FreeVulkanSemaphore(&pipeline.RenderFinishedSemaphores[i]); err != nil {
			return err
		}
		if err := semaphore.FreeVulkanSemaphore(&pipeline.ImageAvailableSemaphores[i]); err != nil {
			return err
		}
		if err := command.FreeVulkanCommandBuffer(&pipeline.CommandBuffers[i]); err != nil {
			return err
		}
	}
	if err := pass.FreeVulkanRenderPass(&pipeline.RenderPass); err != nil {
		return err
	}
	if err := command.FreeVulkanCommandPool(&pipeline.CommandPool); err != nil {
		return err
	}
	if err := layout.FreeVulkanPipelineLayout(&pipeline.Layout); err != nil {
		return err
	}
	vulkan.DestroyPipeline(pipeline.Device.Handle, pipeline.Handle, nil)
	return nil
}
