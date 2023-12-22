package descriptor

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanDescriptorSet struct {
	Handle  vulkan.DescriptorSet
	Device  *logical.VulkanLogicalDevice
	Options VulkanDescriptorSetOptions
}

func NewVulkanDescriptorSet(device *logical.VulkanLogicalDevice, pool *VulkanDescriptorPool, setLayout *VulkanDescriptorSetLayout) (VulkanDescriptorSet, error) {
	descriptorSet := VulkanDescriptorSet{
		Device:  device,
		Options: NewVulkanDescriptorSetOptions(pool, setLayout),
	}

	var vulkanDescriptorSet vulkan.DescriptorSet
	if res := vulkan.AllocateDescriptorSets(device.Handle, &descriptorSet.Options.AllocateInfo, &vulkanDescriptorSet); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return descriptorSet, vulkan.Error(res)
	}
	descriptorSet.Handle = vulkanDescriptorSet
	logger.DefaultLogger.Debug("allocated new vulkan descriptor set")

	return descriptorSet, nil
}

func FreeVulkanDescriptorSet(set *VulkanDescriptorSet) error {
	vulkan.FreeDescriptorSets(set.Device.Handle, set.Options.AllocateInfo.DescriptorPool, 1, &set.Handle)
	return nil
}
