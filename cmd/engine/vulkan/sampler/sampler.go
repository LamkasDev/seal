package sampler

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSampler struct {
	Handle  vulkan.Sampler
	Device  *logical.VulkanLogicalDevice
	Options VulkanSamplerOptions
}

func NewVulkanSampler(device *logical.VulkanLogicalDevice) (VulkanSampler, error) {
	sampler := VulkanSampler{
		Device:  device,
		Options: NewVulkanSamplerOptions(device),
	}

	var vulkanSampler vulkan.Sampler
	if res := vulkan.CreateSampler(device.Handle, &sampler.Options.CreateInfo, nil, &vulkanSampler); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return sampler, vulkan.Error(res)
	}
	sampler.Handle = vulkanSampler
	logger.DefaultLogger.Debug("created new vulkan sampler")

	return sampler, nil
}

func FreeVulkanSampler(sampler *VulkanSampler) error {
	vulkan.DestroySampler(sampler.Device.Handle, sampler.Handle, nil)
	return nil
}
