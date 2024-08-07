package semaphore

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSemaphore struct {
	Handle  vulkan.Semaphore
	Device  *logical.VulkanLogicalDevice
	Options VulkanSemaphoreOptions
}

func NewVulkanSemaphore(device *logical.VulkanLogicalDevice) (VulkanSemaphore, error) {
	semaphore := VulkanSemaphore{
		Device:  device,
		Options: NewVulkanSemaphoreOptions(),
	}

	var vulkanSemaphore vulkan.Semaphore
	if res := vulkan.CreateSemaphore(device.Handle, &semaphore.Options.CreateInfo, nil, &vulkanSemaphore); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return semaphore, vulkan.Error(res)
	}
	semaphore.Handle = vulkanSemaphore
	logger.DefaultLogger.Debug("created new vulkan semaphore")

	return semaphore, nil
}

func FreeVulkanSemaphore(semaphore *VulkanSemaphore) error {
	vulkan.DestroySemaphore(semaphore.Device.Handle, semaphore.Handle, nil)
	return nil
}
