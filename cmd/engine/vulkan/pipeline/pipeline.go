package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipeline struct {
	Handle     vulkan.Pipeline
	Device     *logical.VulkanLogicalDevice
	Viewport   viewport.VulkanViewport
	Layout     VulkanPipelineLayout
	RenderPass pass.VulkanRenderPass
	Options    VulkanPipelineOptions
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, format vulkan.Format, extent vulkan.Extent2D, container *shader.VulkanShaderContainer) (VulkanPipeline, error) {
	var err error
	pipeline := VulkanPipeline{
		Device: device,
	}

	if pipeline.Viewport, err = viewport.NewVulkanViewport(extent); err != nil {
		return pipeline, err
	}
	if pipeline.Layout, err = NewVulkanPipelineLayout(device); err != nil {
		return pipeline, err
	}
	if pipeline.RenderPass, err = pass.NewVulkanRenderPass(device, format); err != nil {
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
	vulkan.DestroyPipelineLayout(pipeline.Device.Handle, pipeline.Options.CreateInfo.Layout, nil)
	vulkan.DestroyRenderPass(pipeline.Device.Handle, pipeline.Options.CreateInfo.RenderPass, nil)
	return nil
}
