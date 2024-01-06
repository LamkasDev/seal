package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanBuffer struct {
	Handle  vulkan.Buffer
	Memory  vulkan.DeviceMemory
	Device  *logical.VulkanLogicalDevice
	Options VulkanBufferOptions
}

func NewVulkanBuffer(device *logical.VulkanLogicalDevice, data VulkanBufferOptionsData) (VulkanBuffer, error) {
	buffer := VulkanBuffer{
		Device:  device,
		Options: NewVulkanBufferOptions(data),
	}

	var vulkanBuffer vulkan.Buffer
	if res := vulkan.CreateBuffer(device.Handle, &buffer.Options.CreateInfo, nil, &vulkanBuffer); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return buffer, vulkan.Error(res)
	}
	buffer.Handle = vulkanBuffer
	logger.DefaultLogger.Debug("created new vulkan buffer")

	if err := AllocateVulkanBuffer(&buffer); err != nil {
		return buffer, err
	}
	logger.DefaultLogger.Debug("allocated vulkan buffer memory")

	return buffer, nil
}

func AllocateVulkanBuffer(buffer *VulkanBuffer) error {
	var requirements vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(buffer.Device.Handle, buffer.Handle, &requirements)
	requirements.Deref()

	if err := UpdateVulkanBufferOptions(&buffer.Options, buffer.Device, requirements); err != nil {
		return err
	}

	var vulkanBufferMemory vulkan.DeviceMemory
	if res := vulkan.AllocateMemory(buffer.Device.Handle, &buffer.Options.AllocateInfo, nil, &vulkanBufferMemory); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	buffer.Memory = vulkanBufferMemory

	if res := vulkan.BindBufferMemory(buffer.Device.Handle, buffer.Handle, buffer.Memory, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func FreeVulkanBuffer(buffer *VulkanBuffer) error {
	vulkan.DestroyBuffer(buffer.Device.Handle, buffer.Handle, nil)
	vulkan.FreeMemory(buffer.Device.Handle, buffer.Memory, nil)
	return nil
}
