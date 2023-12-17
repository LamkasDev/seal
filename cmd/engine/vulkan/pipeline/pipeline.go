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
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipeline struct {
	Handle                  vulkan.Pipeline
	Device                  *logical.VulkanLogicalDevice
	Viewport                viewport.VulkanViewport
	Layout                  layout.VulkanPipelineLayout
	RenderPass              pass.VulkanRenderPass
	CommandPool             command.VulkanCommandPool
	CommandBuffer           command.VulkanCommandBuffer
	SemaphoreImageAvailable semaphore.VulkanSemaphore
	SemaphoreRenderFinished semaphore.VulkanSemaphore
	InFlightFence           fence.VulkanFence
	Options                 VulkanPipelineOptions
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, container *shader.VulkanShaderContainer) (VulkanPipeline, error) {
	var err error
	pipeline := VulkanPipeline{
		Device: device,
	}

	if pipeline.Viewport, err = viewport.NewVulkanViewport(device.Physical.Capabilities.Surface.ImageExtent); err != nil {
		return pipeline, err
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
	if pipeline.CommandBuffer, err = command.NewVulkanCommandBuffer(device, &pipeline.CommandPool); err != nil {
		return pipeline, err
	}
	if pipeline.SemaphoreImageAvailable, err = semaphore.NewVulkanSemaphore(device); err != nil {
		return pipeline, err
	}
	if pipeline.SemaphoreRenderFinished, err = semaphore.NewVulkanSemaphore(device); err != nil {
		return pipeline, err
	}
	if pipeline.InFlightFence, err = fence.NewVulkanFence(device, vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit)); err != nil {
		return pipeline, err
	}
	if pipeline.Options, err = NewVulkanPipelineOptions(&pipeline.Layout, &pipeline.Viewport, &pipeline.RenderPass, container); err != nil {
		return pipeline, err
	}

	vulkanPipelines := make([]vulkan.Pipeline, 1)
	if res := vulkan.CreateGraphicsPipelines(device.Handle, nil, 1, []vulkan.GraphicsPipelineCreateInfo{pipeline.Options.CreateInfo}, nil, vulkanPipelines); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	pipeline.Handle = vulkanPipelines[0]
	logger.DefaultLogger.Debug("created new vulkan pipeline")

	return pipeline, nil
}

func FreeVulkanPipeline(pipeline *VulkanPipeline) error {
	if err := fence.FreeVulkanFence(&pipeline.InFlightFence); err != nil {
		return err
	}
	if err := semaphore.FreeVulkanSemaphore(&pipeline.SemaphoreRenderFinished); err != nil {
		return err
	}
	if err := semaphore.FreeVulkanSemaphore(&pipeline.SemaphoreImageAvailable); err != nil {
		return err
	}
	if err := pass.FreeVulkanRenderPass(&pipeline.RenderPass); err != nil {
		return err
	}
	if err := command.FreeVulkanCommandBuffer(&pipeline.CommandBuffer); err != nil {
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
