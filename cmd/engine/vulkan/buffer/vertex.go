package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanVertexBuffer struct {
	BufferHandle vulkan.Buffer
	MemoryHandle vulkan.DeviceMemory
	Device       *logical.VulkanLogicalDevice
	Options      VulkanVertexBufferOptions
}

func NewVulkanVertexBuffer(device *logical.VulkanLogicalDevice) (VulkanVertexBuffer, error) {
	vertexBuffer := VulkanVertexBuffer{
		Device:  device,
		Options: NewVulkanVertexBufferOptions(vulkan.DeviceSize(int(vertex.VulkanVertexSize) * len(vertex.DefaultVertices))),
	}

	var vulkanVertexBuffer vulkan.Buffer
	if res := vulkan.CreateBuffer(device.Handle, &vertexBuffer.Options.CreateInfo, nil, &vulkanVertexBuffer); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vertexBuffer, vulkan.Error(res)
	}
	vertexBuffer.BufferHandle = vulkanVertexBuffer
	logger.DefaultLogger.Debug("created new vulkan vertex buffer")

	var requirements vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(device.Handle, vertexBuffer.BufferHandle, &requirements)
	if err := UpdateVulkanVertexBufferOptions(&vertexBuffer.Options, device, requirements); err != nil {
		return vertexBuffer, err
	}

	var vulkanVertexBufferMemory vulkan.DeviceMemory
	if res := vulkan.AllocateMemory(device.Handle, &vertexBuffer.Options.AllocateInfo, nil, &vulkanVertexBufferMemory); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vertexBuffer, vulkan.Error(res)
	}
	vertexBuffer.MemoryHandle = vulkanVertexBufferMemory
	logger.DefaultLogger.Debug("allocated vulkan vertex buffer memory")

	if res := vulkan.BindBufferMemory(device.Handle, vertexBuffer.BufferHandle, vertexBuffer.MemoryHandle, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vertexBuffer, vulkan.Error(res)
	}

	return vertexBuffer, nil
}

func FreeVulkanVertexBuffer(buffer *VulkanVertexBuffer) error {
	vulkan.DestroyBuffer(buffer.Device.Handle, buffer.BufferHandle, nil)
	vulkan.FreeMemory(buffer.Device.Handle, buffer.MemoryHandle, nil)
	return nil
}
