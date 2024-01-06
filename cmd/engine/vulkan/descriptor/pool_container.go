package descriptor

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
)

type VulkanDescriptorPoolContainer struct {
	Device *logical.VulkanLogicalDevice
	Pools  []VulkanDescriptorPool
}

func NewVulkanDescriptorPoolContainer(device *logical.VulkanLogicalDevice) (VulkanDescriptorPoolContainer, error) {
	container := VulkanDescriptorPoolContainer{
		Device: device,
		Pools:  []VulkanDescriptorPool{},
	}
	logger.DefaultLogger.Debug("created new vulkan descriptor pool container")

	return container, nil
}

func CreateVulkanDescriptorPoolWithContainer(container *VulkanDescriptorPoolContainer) (VulkanDescriptorPool, error) {
	descriptorPool, err := NewVulkanDescriptorPool(container.Device)
	if err != nil {
		return descriptorPool, err
	}
	container.Pools = append(container.Pools, descriptorPool)

	return descriptorPool, nil
}

func FreeVulkanDescriptorPoolContainer(container *VulkanDescriptorPoolContainer) error {
	for _, pool := range container.Pools {
		if err := FreeVulkanDescriptorPool(&pool); err != nil {
			return err
		}
	}

	return nil
}
