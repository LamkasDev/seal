package logical

import (
	"fmt"
	"slices"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
	"golang.org/x/exp/maps"
)

type VulkanLogicalDeviceOptions struct {
	Features        []vulkan.PhysicalDeviceFeatures
	QueueCreateInfo map[uint32]vulkan.DeviceQueueCreateInfo
	CreateInfo      vulkan.DeviceCreateInfo
}

func NewVulkanLogicalDeviceOptions(device *physical.VulkanPhysicalDevice) VulkanLogicalDeviceOptions {
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
	if slices.ContainsFunc(device.Capabilities.ExtensionNames, func(s string) bool {
		return s == "VK_EXT_memory_priority" || s == "VK_EXT_pageable_device_local_memory"
	}) {
		options.CreateInfo.PpEnabledExtensionNames = append(options.CreateInfo.PpEnabledExtensionNames, "VK_EXT_memory_priority\x00")
		options.CreateInfo.PpEnabledExtensionNames = append(options.CreateInfo.PpEnabledExtensionNames, "VK_EXT_pageable_device_local_memory\x00")
	}
	options.CreateInfo.EnabledExtensionCount = uint32(len(options.CreateInfo.PpEnabledExtensionNames))
	options.CreateInfo.QueueCreateInfoCount = uint32(len(options.CreateInfo.PQueueCreateInfos))

	return options
}

func AddVulkanLogicalDeviceOptionsQueue(options *VulkanLogicalDeviceOptions, index uint32) {
	options.QueueCreateInfo[index] = vulkan.DeviceQueueCreateInfo{
		SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: index,
		QueueCount:       1,
		PQueuePriorities: []float32{1},
	}
}
