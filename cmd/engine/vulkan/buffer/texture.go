package buffer

import (
	"unsafe"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanTextureBuffer struct {
	StagingBuffer VulkanBuffer
	DeviceBuffer  VulkanBuffer
	Device        *logical.VulkanLogicalDevice
	Options       VulkanTextureBufferOptions
}

func NewVulkanTextureBuffer(device *logical.VulkanLogicalDevice, options VulkanTextureBufferOptions) (VulkanTextureBuffer, error) {
	var err error
	textureBuffer := VulkanTextureBuffer{
		Device:  device,
		Options: options,
	}

	if textureBuffer.StagingBuffer, err = NewVulkanBuffer(device, VulkanBufferOptionsData{GetVulkanTextureBufferOptionsSize(&options), vulkan.BufferUsageFlags(vulkan.BufferUsageTransferSrcBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit)}); err != nil {
		return textureBuffer, err
	}
	if textureBuffer.DeviceBuffer, err = NewVulkanBuffer(device, VulkanBufferOptionsData{GetVulkanTextureBufferOptionsSize(&options), vulkan.BufferUsageFlags(vulkan.BufferUsageTransferDstBit | vulkan.BufferUsageVertexBufferBit | vulkan.BufferUsageIndexBufferBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit)}); err != nil {
		return textureBuffer, err
	}
	if err := CopyVulkanTextureBuffers(&textureBuffer); err != nil {
		return textureBuffer, err
	}

	return textureBuffer, nil
}

func CopyVulkanTextureBuffers(buffer *VulkanTextureBuffer) error {
	var vulkanPixelsData unsafe.Pointer
	if res := vulkan.MapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory, GetVulkanTextureBufferOptionsPixelsOffset(&buffer.Options), GetVulkanTextureBufferOptionsPixelsSize(&buffer.Options), 0, &vulkanPixelsData); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	vulkanPixelsBuffer := unsafe.Slice((*uint32)(vulkanPixelsData), len(buffer.Options.Pixels))
	copy(vulkanPixelsBuffer, buffer.Options.Pixels)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	return nil
}

func FreeVulkanTextureBuffer(buffer *VulkanMeshBuffer) error {
	if err := FreeVulkanBuffer(&buffer.StagingBuffer); err != nil {
		return err
	}

	return FreeVulkanBuffer(&buffer.DeviceBuffer)
}
