package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline_layout"
	sealShader "github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/viewport"
	sealWindow "github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipeline struct {
	Handle     vulkan.Pipeline
	Device     *logical.VulkanLogicalDevice
	Window     *sealWindow.Window
	Layout     *pipeline_layout.VulkanPipelineLayout
	Viewport   *viewport.VulkanViewport
	RenderPass *pass.VulkanRenderPass
	Shader     *sealShader.VulkanShader
	Options    VulkanPipelineOptions
}

func NewVulkanPipeline(device *logical.VulkanLogicalDevice, window *sealWindow.Window, layout *pipeline_layout.VulkanPipelineLayout, viewport *viewport.VulkanViewport, renderPass *pass.VulkanRenderPass, shader *sealShader.VulkanShader) (VulkanPipeline, error) {
	pipeline := VulkanPipeline{
		Device:     device,
		Window:     window,
		Layout:     layout,
		Viewport:   viewport,
		RenderPass: renderPass,
		Shader:     shader,
		Options:    NewVulkanPipelineOptions(layout, viewport, renderPass, shader),
	}

	vulkanPipelines := make([]vulkan.Pipeline, 1)
	if res := vulkan.CreateGraphicsPipelines(device.Handle, nil, 1, []vulkan.GraphicsPipelineCreateInfo{pipeline.Options.CreateInfo}, nil, vulkanPipelines); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pipeline, vulkan.Error(res)
	}
	pipeline.Handle = vulkanPipelines[0]
	logger.DefaultLogger.Debug("created new vulkan pipeline")

	return pipeline, nil
}

func FreeVulkanPipeline(pipeline *VulkanPipeline) error {
	vulkan.DestroyPipeline(pipeline.Device.Handle, pipeline.Handle, nil)
	return nil
}
