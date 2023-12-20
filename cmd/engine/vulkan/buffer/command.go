package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandBuffer struct {
	Handle  vulkan.CommandBuffer
	Device  *logical.VulkanLogicalDevice
	Pool    *command.VulkanCommandPool
	Options VulkanCommandBufferOptions
}

func NewVulkanCommandBuffer(device *logical.VulkanLogicalDevice, pool *command.VulkanCommandPool) (VulkanCommandBuffer, error) {
	commandBuffer := VulkanCommandBuffer{
		Device:  device,
		Pool:    pool,
		Options: NewVulkanCommandBufferOptions(pool),
	}

	vulkanCommandBuffers := make([]vulkan.CommandBuffer, 1)
	if res := vulkan.AllocateCommandBuffers(device.Handle, &commandBuffer.Options.AllocateInfo, vulkanCommandBuffers); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return commandBuffer, vulkan.Error(res)
	}
	commandBuffer.Handle = vulkanCommandBuffers[0]
	logger.DefaultLogger.Debug("created new vulkan command buffer")

	return commandBuffer, nil
}

func BeginVulkanCommandBuffer(buffer *VulkanCommandBuffer) error {
	if res := vulkan.ResetCommandBuffer(buffer.Handle, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	beginInfo := vulkan.CommandBufferBeginInfo{
		SType: vulkan.StructureTypeCommandBufferBeginInfo,
		Flags: vulkan.CommandBufferUsageFlags(vulkan.CommandBufferUsageOneTimeSubmitBit),
	}
	if res := vulkan.BeginCommandBuffer(buffer.Handle, &beginInfo); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func EndVulkanCommandBuffer(buffer *VulkanCommandBuffer) error {
	if res := vulkan.EndCommandBuffer(buffer.Handle); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func FreeVulkanCommandBuffer(buffer *VulkanCommandBuffer) error {
	vulkan.FreeCommandBuffers(buffer.Device.Handle, buffer.Pool.Handle, 1, []vulkan.CommandBuffer{buffer.Handle})
	return nil
}
