package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayout struct {
	Handle  vulkan.PipelineLayout
	Options VulkanPipelineLayoutOptions
}

func NewVulkanPipelineLayout(device *logical.VulkanLogicalDevice) (VulkanPipelineLayout, error) {
	var err error
	pipelineLayout := VulkanPipelineLayout{}

	if pipelineLayout.Options, err = NewVulkanPipelineLayoutOptions(); err != nil {
		return pipelineLayout, err
	}
	logger.DefaultLogger.Debug("created new vulkan pipeline layout options")

	var vulkanPipelineLayout vulkan.PipelineLayout
	if res := vulkan.CreatePipelineLayout(device.Handle, &pipelineLayout.Options.CreateInfo, nil, &vulkanPipelineLayout); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	pipelineLayout.Handle = vulkanPipelineLayout
	logger.DefaultLogger.Debug("created new vulkan pipeline layout")

	return pipelineLayout, nil
}
