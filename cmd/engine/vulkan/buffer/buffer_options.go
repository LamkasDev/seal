package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
)

type VulkanBufferOptionsData struct {
	Size        vulkan.DeviceSize
	Usage       vulkan.BufferUsageFlags
	SharingMode vulkan.SharingMode
	Flags       vulkan.MemoryPropertyFlags
}

type VulkanBufferOptions struct {
	Data         VulkanBufferOptionsData
	CreateInfo   vulkan.BufferCreateInfo
	AllocateInfo vulkan.MemoryAllocateInfo
}

func NewVulkanBufferOptions(data VulkanBufferOptionsData) VulkanBufferOptions {
	options := VulkanBufferOptions{
		Data: data,
		CreateInfo: vulkan.BufferCreateInfo{
			SType:       vulkan.StructureTypeBufferCreateInfo,
			Size:        data.Size,
			Usage:       data.Usage,
			SharingMode: data.SharingMode,
		},
	}

	return options
}

func UpdateVulkanBufferOptions(options *VulkanBufferOptions, device *logical.VulkanLogicalDevice, requirements vulkan.MemoryRequirements) error {
	memoryTypeIndex, err := physical.GetVulkanPhysicalDeviceMemoryTypeIndex(&device.Physical.Capabilities, requirements.MemoryTypeBits, options.Data.Flags)
	if err != nil {
		return err
	}
	options.AllocateInfo = vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	return nil
}
