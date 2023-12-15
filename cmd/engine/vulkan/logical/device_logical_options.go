package logical

import (
	"fmt"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
	"golang.org/x/exp/maps"
)

type VulkanLogicalDeviceOptions struct {
	Features        []vulkan.PhysicalDeviceFeatures
	QueueCreateInfo map[uint32]vulkan.DeviceQueueCreateInfo
	CreateInfo      vulkan.DeviceCreateInfo
}

func NewVulkanLogicalDeviceOptions(device *physical.VulkanPhysicalDevice) (VulkanLogicalDeviceOptions, error) {
	options := VulkanLogicalDeviceOptions{
		Features: []vulkan.PhysicalDeviceFeatures{
			{},
		},
		QueueCreateInfo: map[uint32]vulkan.DeviceQueueCreateInfo{},
	}
	AddVulkanLogicalDeviceOptionsQueue(&options, uint32(device.Capabilities.Queue.GraphicsIndex))
	AddVulkanLogicalDeviceOptionsQueue(&options, uint32(device.Capabilities.Queue.PresentationIndex))
	options.CreateInfo = vulkan.DeviceCreateInfo{
		SType:            vulkan.StructureTypeDeviceCreateInfo,
		PEnabledFeatures: options.Features,
		PpEnabledExtensionNames: []string{
			fmt.Sprintf("%s\x00", vulkan.KhrSwapchainExtensionName),
		},
		PQueueCreateInfos: maps.Values(options.QueueCreateInfo),
	}
	options.CreateInfo.EnabledExtensionCount = uint32(len(options.CreateInfo.PpEnabledExtensionNames))
	options.CreateInfo.QueueCreateInfoCount = uint32(len(options.CreateInfo.PQueueCreateInfos))

	return options, nil
}

func AddVulkanLogicalDeviceOptionsQueue(options *VulkanLogicalDeviceOptions, index uint32) {
	options.QueueCreateInfo[index] = vulkan.DeviceQueueCreateInfo{
		SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: index,
		QueueCount:       1,
		PQueuePriorities: []float32{1},
	}
}
