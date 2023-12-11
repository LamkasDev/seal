package vulkan

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanLogicalDeviceOptions struct {
	QueueCreateInfo []vulkan.DeviceQueueCreateInfo
	Features        []vulkan.PhysicalDeviceFeatures
	CreateInfo      vulkan.DeviceCreateInfo
}

func NewVulkanLogicalDeviceOptions(device *VulkanPhysicalDevice) (VulkanLogicalDeviceOptions, error) {
	options := VulkanLogicalDeviceOptions{}
	options.QueueCreateInfo = []vulkan.DeviceQueueCreateInfo{
		{
			SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: uint32(device.QueueFamilyGraphicsIndex),
			QueueCount:       1,
			PQueuePriorities: []float32{1},
		},
	}
	options.Features = []vulkan.PhysicalDeviceFeatures{
		{},
	}
	options.CreateInfo = vulkan.DeviceCreateInfo{
		SType:             vulkan.StructureTypeDeviceCreateInfo,
		PQueueCreateInfos: options.QueueCreateInfo,
		PEnabledFeatures:  options.Features,
	}
	options.CreateInfo.QueueCreateInfoCount = uint32(len(options.CreateInfo.PQueueCreateInfos))

	return options, nil
}
