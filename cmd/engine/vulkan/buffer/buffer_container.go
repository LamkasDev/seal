package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
)

type VulkanBufferContainer struct {
	Device  *logical.VulkanLogicalDevice
	Buffers []VulkanBuffer
}

func NewVulkanBufferContainer(device *logical.VulkanLogicalDevice) (VulkanBufferContainer, error) {
	container := VulkanBufferContainer{
		Device:  device,
		Buffers: []VulkanBuffer{},
	}
	logger.DefaultLogger.Debug("created new vulkan buffer container")

	return container, nil
}

func CreateVulkanBufferWithContainer(container *VulkanBufferContainer, data VulkanBufferOptionsData) (VulkanBuffer, error) {
	buffer, err := NewVulkanBuffer(container.Device, data)
	if err != nil {
		return buffer, err
	}
	container.Buffers = append(container.Buffers, buffer)

	return buffer, nil
}

func FreeVulkanBufferContainer(container *VulkanBufferContainer) error {
	for _, buffer := range container.Buffers {
		if err := FreeVulkanBuffer(&buffer); err != nil {
			return err
		}
	}

	return nil
}
