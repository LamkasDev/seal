package buffer

import (
	"unsafe"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/vulkan-go/vulkan"
)

type VulkanVertexBuffer struct {
	StagingBuffer VulkanBuffer
	DeviceBuffer  VulkanBuffer
	Device        *logical.VulkanLogicalDevice
}

func NewVulkanVertexBuffer(device *logical.VulkanLogicalDevice) (VulkanVertexBuffer, error) {
	var err error
	vertexBuffer := VulkanVertexBuffer{
		Device: device,
	}

	bufferSize := vulkan.DeviceSize(int(vertex.VulkanVertexSize) * len(vertex.DefaultVertices))
	if vertexBuffer.StagingBuffer, err = NewVulkanBuffer(device, bufferSize, vulkan.BufferUsageFlags(vulkan.BufferUsageTransferSrcBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit)); err != nil {
		return vertexBuffer, err
	}
	if vertexBuffer.DeviceBuffer, err = NewVulkanBuffer(device, bufferSize, vulkan.BufferUsageFlags(vulkan.BufferUsageTransferDstBit|vulkan.BufferUsageVertexBufferBit), vulkan.SharingModeExclusive, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit)); err != nil {
		return vertexBuffer, err
	}
	if err := SetVulkanVertexBuffer(&vertexBuffer, vertex.DefaultVertices); err != nil {
		return vertexBuffer, err
	}

	return vertexBuffer, nil
}

func SetVulkanVertexBuffer(buffer *VulkanVertexBuffer, vertices []vertex.VulkanVertex) error {
	var vulkanVertexVerticesData unsafe.Pointer
	vulkan.MapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory, 0, buffer.StagingBuffer.Options.CreateInfo.Size, 0, &vulkanVertexVerticesData)
	vulkanVertexBufferVertices := unsafe.Slice((*vertex.VulkanVertex)(vulkanVertexVerticesData), len(vertices))
	copy(vulkanVertexBufferVertices, vertex.DefaultVertices)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	return nil
}

func FreeVulkanVertexBuffer(buffer *VulkanVertexBuffer) error {
	if err := FreeVulkanBuffer(&buffer.StagingBuffer); err != nil {
		return err
	}

	return FreeVulkanBuffer(&buffer.DeviceBuffer)
}
