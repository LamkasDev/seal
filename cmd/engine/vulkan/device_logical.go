package vulkan

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanLogicalDevice struct {
	Handle   vulkan.Device
	Physical *VulkanPhysicalDevice
	Options  VulkanLogicalDeviceOptions
	Queues   map[uint32]vulkan.Queue
}

func NewVulkanLogicalDevice(physicalDevice *VulkanPhysicalDevice) (VulkanLogicalDevice, error) {
	var err error
	device := VulkanLogicalDevice{
		Physical: physicalDevice,
		Queues:   map[uint32]vulkan.Queue{},
	}

	if device.Options, err = NewVulkanLogicalDeviceOptions(physicalDevice); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan logical device options")

	var vulkanDevice vulkan.Device
	vulkan.CreateDevice(physicalDevice.Handle, &device.Options.CreateInfo, nil, &vulkanDevice)
	device.Handle = vulkanDevice
	logger.DefaultLogger.Info("created new vulkan logical device")

	for k := range device.Options.QueueCreateInfo {
		var vulkanQueue vulkan.Queue
		vulkan.GetDeviceQueue(device.Handle, k, 0, &vulkanQueue)
		device.Queues[k] = vulkanQueue
	}

	return device, nil
}

func FreeVulkanLogicalDevice(device *VulkanLogicalDevice) error {
	vulkan.DestroyDevice(device.Handle, nil)
	return nil
}
