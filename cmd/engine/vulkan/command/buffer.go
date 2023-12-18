package command

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanCommandBuffer struct {
	Handle  vulkan.CommandBuffer
	Device  *logical.VulkanLogicalDevice
	Pool    *VulkanCommandPool
	Options VulkanCommandBufferOptions
}

func NewVulkanCommandBuffer(device *logical.VulkanLogicalDevice, pool *VulkanCommandPool) (VulkanCommandBuffer, error) {
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

func FreeVulkanCommandBuffer(buffer *VulkanCommandBuffer) error {
	vulkan.FreeCommandBuffers(buffer.Device.Handle, buffer.Pool.Handle, 1, []vulkan.CommandBuffer{buffer.Handle})
	return nil
}
