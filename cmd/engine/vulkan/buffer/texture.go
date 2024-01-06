package buffer

import (
	"image"
	"unsafe"

	sealImage "github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanTextureBuffer struct {
	StagingBuffer VulkanBuffer
	Image         *sealImage.VulkanImage
	Device        *logical.VulkanLogicalDevice
	Options       VulkanTextureBufferOptions
}

func NewVulkanTextureBuffer(device *logical.VulkanLogicalDevice, source *image.RGBA, image *sealImage.VulkanImage) (VulkanTextureBuffer, error) {
	var err error
	textureBuffer := VulkanTextureBuffer{
		Image:  image,
		Device: device,
		Options: VulkanTextureBufferOptions{
			Pixels: source.Pix,
		},
	}

	if textureBuffer.StagingBuffer, err = NewVulkanBuffer(device, VulkanBufferOptionsData{
		Size:        image.Options.AllocateInfo.AllocationSize,
		Usage:       vulkan.BufferUsageFlags(vulkan.BufferUsageTransferSrcBit),
		SharingMode: vulkan.SharingModeExclusive,
		Flags:       vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit),
	}); err != nil {
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
	vulkanPixelsBuffer := unsafe.Slice((*uint8)(vulkanPixelsData), len(buffer.Options.Pixels))
	copy(vulkanPixelsBuffer, buffer.Options.Pixels)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	return nil
}

func FreeVulkanTextureBuffer(buffer *VulkanTextureBuffer) error {
	return FreeVulkanBuffer(&buffer.StagingBuffer)
}
