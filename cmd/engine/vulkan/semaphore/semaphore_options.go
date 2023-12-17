package semaphore

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanSemaphoreOptions struct {
	CreateInfo vulkan.SemaphoreCreateInfo
}

func NewVulkanSemaphoreOptions() (VulkanSemaphoreOptions, error) {
	options := VulkanSemaphoreOptions{
		CreateInfo: vulkan.SemaphoreCreateInfo{
			SType: vulkan.StructureTypeSemaphoreCreateInfo,
		},
	}

	return options, nil
}
