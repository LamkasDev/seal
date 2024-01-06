package buffer

import (
	"unsafe"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanMeshBuffer struct {
	StagingBuffer VulkanBuffer
	DeviceBuffer  VulkanBuffer
	Device        *logical.VulkanLogicalDevice
	Options       VulkanMeshBufferOptions
}

func NewVulkanMeshBuffer(device *logical.VulkanLogicalDevice, options VulkanMeshBufferOptions) (VulkanMeshBuffer, error) {
	var err error
	meshBuffer := VulkanMeshBuffer{
		Device:  device,
		Options: options,
	}

	if meshBuffer.StagingBuffer, err = NewVulkanBuffer(device, VulkanBufferOptionsData{GetVulkanMeshBufferOptionsSize(&options), vulkan.BufferUsageFlags(vulkan.BufferUsageTransferSrcBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit)}); err != nil {
		return meshBuffer, err
	}
	if meshBuffer.DeviceBuffer, err = NewVulkanBuffer(device, VulkanBufferOptionsData{GetVulkanMeshBufferOptionsSize(&options), vulkan.BufferUsageFlags(vulkan.BufferUsageTransferDstBit | vulkan.BufferUsageVertexBufferBit | vulkan.BufferUsageIndexBufferBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit)}); err != nil {
		return meshBuffer, err
	}
	if err := CopyVulkanMeshBuffers(&meshBuffer); err != nil {
		return meshBuffer, err
	}

	return meshBuffer, nil
}

func CopyVulkanMeshBuffers(buffer *VulkanMeshBuffer) error {
	var vulkanVerticesData unsafe.Pointer
	if res := vulkan.MapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory, GetVulkanMeshBufferOptionsVerticesOffset(&buffer.Options), GetVulkanMeshBufferOptionsVerticesSize(&buffer.Options), 0, &vulkanVerticesData); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	vulkanVerticesBuffer := unsafe.Slice((*vertex.VulkanVertex)(vulkanVerticesData), len(buffer.Options.Vertices))
	copy(vulkanVerticesBuffer, buffer.Options.Vertices)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	var vulkanIndexData unsafe.Pointer
	if res := vulkan.MapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory, GetVulkanMeshBufferOptionsIndicesOffset(&buffer.Options), GetVulkanMeshBufferOptionsIndicesSize(&buffer.Options), 0, &vulkanIndexData); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	vulkanIndexBuffer := unsafe.Slice((*uint16)(vulkanIndexData), len(buffer.Options.Indices))
	copy(vulkanIndexBuffer, buffer.Options.Indices)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	return nil
}

func FreeVulkanMeshBuffer(buffer *VulkanMeshBuffer) error {
	if err := FreeVulkanBuffer(&buffer.StagingBuffer); err != nil {
		return err
	}

	return FreeVulkanBuffer(&buffer.DeviceBuffer)
}
