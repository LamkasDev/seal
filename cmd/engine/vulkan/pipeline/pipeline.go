package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/layout"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

const MaxFramesInFlight = 2

type VulkanPipeline struct {
	Handle       vulkan.Pipeline
	Device       *logical.VulkanLogicalDevice
	Container    *shader.VulkanShaderContainer
	Window       *window.Window
	Viewport     viewport.VulkanViewport
	Layout       layout.VulkanPipelineLayout
	RenderPass   pass.VulkanRenderPass
	Syncer       VulkanPipelineSyncer
	Commander    VulkanPipelineCommander
	VertexBuffer buffer.VulkanVertexBuffer
	CurrentFrame uint32
	Options      VulkanPipelineOptions
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, container *shader.VulkanShaderContainer, cwindow *window.Window) (VulkanPipeline, error) {
	var err error
	pipeline := VulkanPipeline{
		Device:    device,
		Container: container,
		Window:    cwindow,
		Viewport:  viewport.NewVulkanViewport(cwindow.Data.Extent),
	}

	if pipeline.Layout, err = layout.NewVulkanPipelineLayout(device); err != nil {
		return pipeline, err
	}
	if pipeline.RenderPass, err = pass.NewVulkanRenderPass(device, device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format); err != nil {
		return pipeline, err
	}
	if pipeline.Syncer, err = NewVulkanPipelineSyncer(device); err != nil {
		return pipeline, err
	}
	if pipeline.Commander, err = NewVulkanPipelineCommander(device); err != nil {
		return pipeline, err
	}
	if pipeline.VertexBuffer, err = buffer.NewVulkanVertexBuffer(device); err != nil {
		return pipeline, err
	}
	pipeline.Options = NewVulkanPipelineOptions(&pipeline.Layout, &pipeline.Viewport, &pipeline.RenderPass, container)

	vulkanPipelines := make([]vulkan.Pipeline, 1)
	if res := vulkan.CreateGraphicsPipelines(device.Handle, nil, 1, []vulkan.GraphicsPipelineCreateInfo{pipeline.Options.CreateInfo}, nil, vulkanPipelines); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pipeline, vulkan.Error(res)
	}
	pipeline.Handle = vulkanPipelines[0]
	logger.DefaultLogger.Debug("created new vulkan pipeline")

	// Copy default vertices to vertex buffer
	if err := PushVulkanPipelineBuffers(&pipeline); err != nil {
		return pipeline, err
	}

	return pipeline, nil
}

func PushVulkanPipelineBuffers(pipeline *VulkanPipeline) error {
	if err := buffer.BeginVulkanCommandBuffer(&pipeline.Commander.StagingBuffer); err != nil {
		return err
	}

	bufferCopy := vulkan.BufferCopy{
		Size: pipeline.VertexBuffer.StagingBuffer.Options.CreateInfo.Size,
	}
	vulkan.CmdCopyBuffer(pipeline.Commander.StagingBuffer.Handle, pipeline.VertexBuffer.StagingBuffer.Handle, pipeline.VertexBuffer.DeviceBuffer.Handle, 1, []vulkan.BufferCopy{bufferCopy})

	if err := buffer.EndVulkanCommandBuffer(&pipeline.Commander.StagingBuffer); err != nil {
		return err
	}

	submitInfo := vulkan.SubmitInfo{
		SType:              vulkan.StructureTypeSubmitInfo,
		CommandBufferCount: 1,
		PCommandBuffers:    []vulkan.CommandBuffer{pipeline.Commander.StagingBuffer.Handle},
	}
	if res := vulkan.QueueSubmit(pipeline.Device.Queues[uint32(pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)], 1, []vulkan.SubmitInfo{submitInfo}, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	if res := vulkan.QueueWaitIdle(pipeline.Device.Queues[uint32(pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex)]); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func ResizeVulkanPipeline(pipeline *VulkanPipeline) error {
	pipeline.Viewport = viewport.NewVulkanViewport(pipeline.Window.Data.Extent)
	return nil
}

func FreeVulkanPipeline(pipeline *VulkanPipeline) error {
	if err := buffer.FreeVulkanVertexBuffer(&pipeline.VertexBuffer); err != nil {
		return err
	}
	if err := FreeVulkanPipelineSyncer(&pipeline.Syncer); err != nil {
		return err
	}
	if err := FreeVulkanPipelineCommander(&pipeline.Commander); err != nil {
		return err
	}
	if err := pass.FreeVulkanRenderPass(&pipeline.RenderPass); err != nil {
		return err
	}
	if err := layout.FreeVulkanPipelineLayout(&pipeline.Layout); err != nil {
		return err
	}
	vulkan.DestroyPipeline(pipeline.Device.Handle, pipeline.Handle, nil)
	return nil
}
