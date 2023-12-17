package logical

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanLogicalDevice struct {
	Handle   vulkan.Device
	Physical *physical.VulkanPhysicalDevice
	Options  VulkanLogicalDeviceOptions
	Queues   map[uint32]vulkan.Queue
}

func NewVulkanLogicalDevice(physicalDevice *physical.VulkanPhysicalDevice) (VulkanLogicalDevice, error) {
	var err error
	device := VulkanLogicalDevice{
		Physical: physicalDevice,
		Queues:   map[uint32]vulkan.Queue{},
	}

	if device.Options, err = NewVulkanLogicalDeviceOptions(physicalDevice); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	var vulkanDevice vulkan.Device
	if res := vulkan.CreateDevice(physicalDevice.Handle, &device.Options.CreateInfo, nil, &vulkanDevice); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	device.Handle = vulkanDevice
	logger.DefaultLogger.Debug("created new vulkan logical device")

	for k := range device.Options.QueueCreateInfo {
		var vulkanQueue vulkan.Queue
		vulkan.GetDeviceQueue(device.Handle, k, 0, &vulkanQueue)
		device.Queues[k] = vulkanQueue
	}
	logger.DefaultLogger.Debug("retrieved vulkan logical device queues")

	return device, nil
}

func FreeVulkanLogicalDevice(device *VulkanLogicalDevice) error {
	vulkan.DestroyDevice(device.Handle, nil)
	return nil
}
