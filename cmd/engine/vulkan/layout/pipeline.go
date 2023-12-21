package layout

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayout struct {
	Handle        vulkan.PipelineLayout
	DescriptorSet VulkanDescriptorSetLayout
	Device        *logical.VulkanLogicalDevice
	Options       VulkanPipelineLayoutOptions
}

func NewVulkanPipelineLayout(device *logical.VulkanLogicalDevice) (VulkanPipelineLayout, error) {
	var err error
	pipelineLayout := VulkanPipelineLayout{
		Device: device,
	}

	if pipelineLayout.DescriptorSet, err = NewVulkanDescriptorSetLayout(device); err != nil {
		return pipelineLayout, err
	}
	pipelineLayout.Options = NewVulkanPipelineLayoutOptions(&pipelineLayout.DescriptorSet)

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
	if err := FreeVulkanDescriptorSetLayout(&layout.DescriptorSet); err != nil {
		return err
	}

	vulkan.DestroyPipelineLayout(layout.Device.Handle, layout.Handle, nil)
	return nil
}
