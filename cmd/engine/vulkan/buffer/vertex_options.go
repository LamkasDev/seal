package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
)

type VulkanVertexBufferOptions struct {
	CreateInfo   vulkan.BufferCreateInfo
	AllocateInfo vulkan.MemoryAllocateInfo
}

func NewVulkanVertexBufferOptions(size vulkan.DeviceSize) VulkanVertexBufferOptions {
	options := VulkanVertexBufferOptions{
		CreateInfo: vulkan.BufferCreateInfo{
			SType:       vulkan.StructureTypeBufferCreateInfo,
			Size:        size,
			Usage:       vulkan.BufferUsageFlags(vulkan.BufferUsageVertexBufferBit),
			SharingMode: vulkan.SharingModeExclusive,
		},
	}

	return options
}

func UpdateVulkanVertexBufferOptions(options *VulkanVertexBufferOptions, device *logical.VulkanLogicalDevice, requirements vulkan.MemoryRequirements) error {
	memoryTypeIndex, err := physical.GetVulkanPhysicalDeviceMemoryTypeIndex(&device.Physical.Capabilities, requirements.MemoryTypeBits, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit))
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
