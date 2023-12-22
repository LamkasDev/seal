package descriptor

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorSetLayout struct {
	Handle  vulkan.DescriptorSetLayout
	Device  *logical.VulkanLogicalDevice
	Options VulkanDescriptorSetLayoutOptions
}

func NewVulkanDescriptorSetLayout(device *logical.VulkanLogicalDevice) (VulkanDescriptorSetLayout, error) {
	descriptorSetLayout := VulkanDescriptorSetLayout{
		Device:  device,
		Options: NewVulkanDescriptorSetLayoutOptions(),
	}

	var vulkanDescriptorSetLayout vulkan.DescriptorSetLayout
	if res := vulkan.CreateDescriptorSetLayout(device.Handle, &descriptorSetLayout.Options.CreateInfo, nil, &vulkanDescriptorSetLayout); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return descriptorSetLayout, vulkan.Error(res)
	}
	descriptorSetLayout.Handle = vulkanDescriptorSetLayout
	logger.DefaultLogger.Debug("created new vulkan descriptor set layout")

	return descriptorSetLayout, nil
}

func FreeVulkanDescriptorSetLayout(layout *VulkanDescriptorSetLayout) error {
	vulkan.DestroyDescriptorSetLayout(layout.Device.Handle, layout.Handle, nil)
	return nil
}
