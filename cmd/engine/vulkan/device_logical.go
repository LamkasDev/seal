package vulkan

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanLogicalDevice struct {
	Handle   vulkan.Device
	Physical *VulkanPhysicalDevice
	Options  VulkanLogicalDeviceOptions
	Queue    vulkan.Queue
}

func NewVulkanLogicalDevice(physicalDevice *VulkanPhysicalDevice) (VulkanLogicalDevice, error) {
	var err error
	device := VulkanLogicalDevice{
		Physical: physicalDevice,
	}

	if device.Options, err = NewVulkanLogicalDeviceOptions(physicalDevice); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan logical device options")

	var vulkanDevice vulkan.Device
	vulkan.CreateDevice(physicalDevice.Handle, &device.Options.CreateInfo, nil, &vulkanDevice)
	device.Handle = vulkanDevice
	logger.DefaultLogger.Info("created new vulkan logical device")

	var vulkanQueue vulkan.Queue
	vulkan.GetDeviceQueue(device.Handle, uint32(device.Physical.QueueFamilyGraphicsIndex), 0, &vulkanQueue)
	device.Queue = vulkanQueue

	return device, nil
}

func ProbeVulkanLogicalDevice(device *VulkanLogicalDevice, surface *vulkan.Surface) {
	ProbeVulkanPhysicalDevice(device.Physical, surface)
}

func FreeVulkanLogicalDevice(device *VulkanLogicalDevice) error {
	vulkan.DestroyDevice(device.Handle, nil)
	return nil
}
