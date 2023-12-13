package vulkan

import (
	"errors"
	"slices"

	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceDevices struct {
	PhysicalDevices []VulkanPhysicalDevice
	LogicalDevice   VulkanLogicalDevice
}

func NewVulkanInstanceDevices(instance *VulkanInstance, window *glfw.Window, surface *vulkan.Surface) (VulkanInstanceDevices, error) {
	var err error
	devices := VulkanInstanceDevices{}

	var rawDeviceCount uint32
	vulkan.EnumeratePhysicalDevices(instance.Handle, &rawDeviceCount, nil)
	rawDevices := make([]vulkan.PhysicalDevice, rawDeviceCount)
	vulkan.EnumeratePhysicalDevices(instance.Handle, &rawDeviceCount, rawDevices)

	devices.PhysicalDevices = []VulkanPhysicalDevice{}
	for i := 0; i < len(rawDevices); i++ {
		var device VulkanPhysicalDevice
		if device, err = NewVulkanPhysicalDevice(rawDevices[i], window, surface); err != nil {
			logger.DefaultLogger.Warnf("failed to create a new vulkan device")
		}
		if IsVulkanPhysicalDeviceSupported(&device) {
			devices.PhysicalDevices = append(devices.PhysicalDevices, device)
		}
	}
	slices.SortFunc(devices.PhysicalDevices, func(a, b VulkanPhysicalDevice) int {
		return CompareVulkanPhysicalDevice(&b) - CompareVulkanPhysicalDevice(&a)
	})
	logger.DefaultLogger.Infof("found %d/%d suitable physical devices", len(devices.PhysicalDevices), len(rawDevices))
	if len(devices.PhysicalDevices) == 0 {
		return devices, errors.New("no suitable physical device found")
	}

	// this will work, as long as the array doesn't relocate :)
	if devices.LogicalDevice, err = NewVulkanLogicalDevice(&devices.PhysicalDevices[0]); err != nil {
		return devices, err
	}

	return devices, nil
}

func FreeVulkanInstanceDevices(devices *VulkanInstanceDevices) error {
	return FreeVulkanLogicalDevice(&devices.LogicalDevice)
}
