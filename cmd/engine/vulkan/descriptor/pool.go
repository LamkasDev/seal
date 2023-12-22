package descriptor

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorPool struct {
	Handle  vulkan.DescriptorPool
	Device  *logical.VulkanLogicalDevice
	Options VulkanDescriptorPoolOptions
}

func NewVulkanDescriptorPool(device *logical.VulkanLogicalDevice) (VulkanDescriptorPool, error) {
	descriptorPool := VulkanDescriptorPool{
		Device:  device,
		Options: NewVulkanDescriptorPoolOptions(),
	}

	var vulkanDescriptorPool vulkan.DescriptorPool
	if res := vulkan.CreateDescriptorPool(device.Handle, &descriptorPool.Options.CreateInfo, nil, &vulkanDescriptorPool); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return descriptorPool, vulkan.Error(res)
	}
	descriptorPool.Handle = vulkanDescriptorPool
	logger.DefaultLogger.Debug("created new vulkan descriptor pool")

	return descriptorPool, nil
}

func ResetVulkanDescriptorPool(pool *VulkanDescriptorPool) error {
	if res := vulkan.ResetDescriptorPool(pool.Device.Handle, pool.Handle, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func FreeVulkanDescriptorPool(pool *VulkanDescriptorPool) error {
	vulkan.DestroyDescriptorPool(pool.Device.Handle, pool.Handle, nil)
	return nil
}
