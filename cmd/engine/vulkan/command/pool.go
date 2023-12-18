package command

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandPool struct {
	Handle  vulkan.CommandPool
	Device  *logical.VulkanLogicalDevice
	Options VulkanCommandPoolOptions
}

func NewVulkanCommandPool(device *logical.VulkanLogicalDevice, queueFamilyIndex uint32) (VulkanCommandPool, error) {
	commandPool := VulkanCommandPool{
		Device:  device,
		Options: NewVulkanCommandPoolOptions(queueFamilyIndex),
	}

	var vulkanCommandPool vulkan.CommandPool
	if res := vulkan.CreateCommandPool(device.Handle, &commandPool.Options.CreateInfo, nil, &vulkanCommandPool); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return commandPool, vulkan.Error(res)
	}
	commandPool.Handle = vulkanCommandPool
	logger.DefaultLogger.Debug("created new vulkan command pool")

	return commandPool, nil
}

func FreeVulkanCommandPool(pool *VulkanCommandPool) error {
	vulkan.DestroyCommandPool(pool.Device.Handle, pool.Handle, nil)
	return nil
}
