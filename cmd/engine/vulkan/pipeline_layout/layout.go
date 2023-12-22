package pipeline_layout

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/descriptor"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayout struct {
	Handle              vulkan.PipelineLayout
	DescriptorSetLayout descriptor.VulkanDescriptorSetLayout
	Device              *logical.VulkanLogicalDevice
	Options             VulkanPipelineLayoutOptions
}

func NewVulkanPipelineLayout(device *logical.VulkanLogicalDevice) (VulkanPipelineLayout, error) {
	var err error
	pipelineLayout := VulkanPipelineLayout{
		Device: device,
	}

	if pipelineLayout.DescriptorSetLayout, err = descriptor.NewVulkanDescriptorSetLayout(device); err != nil {
		return pipelineLayout, err
	}
	pipelineLayout.Options = NewVulkanPipelineLayoutOptions(&pipelineLayout.DescriptorSetLayout)

	var vulkanPipelineLayout vulkan.PipelineLayout
	if res := vulkan.CreatePipelineLayout(device.Handle, &pipelineLayout.Options.CreateInfo, nil, &vulkanPipelineLayout); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pipelineLayout, vulkan.Error(res)
	}
	pipelineLayout.Handle = vulkanPipelineLayout
	logger.DefaultLogger.Debug("created new vulkan pipeline layout")

	return pipelineLayout, nil
}

func FreeVulkanPipelineLayout(layout *VulkanPipelineLayout) error {
	vulkan.DestroyPipelineLayout(layout.Device.Handle, layout.Handle, nil)
	if err := descriptor.FreeVulkanDescriptorSetLayout(&layout.DescriptorSetLayout); err != nil {
		return err
	}

	return nil
}
