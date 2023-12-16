package device

import (
	"errors"
	"slices"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceDevices struct {
	PhysicalDevices []physical.VulkanPhysicalDevice
	LogicalDevice   logical.VulkanLogicalDevice
}

func NewVulkanInstanceDevices(instance vulkan.Instance, window *glfw.Window, surface *vulkan.Surface) (VulkanInstanceDevices, error) {
	var err error
	devices := VulkanInstanceDevices{}

	var rawDeviceCount uint32
	if res := vulkan.EnumeratePhysicalDevices(instance, &rawDeviceCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	rawDevices := make([]vulkan.PhysicalDevice, rawDeviceCount)
	if res := vulkan.EnumeratePhysicalDevices(instance, &rawDeviceCount, rawDevices); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	devices.PhysicalDevices = []physical.VulkanPhysicalDevice{}
	for i := 0; i < len(rawDevices); i++ {
		var device physical.VulkanPhysicalDevice
		if device, err = physical.NewVulkanPhysicalDevice(rawDevices[i], window, surface); err != nil {
			logger.DefaultLogger.Warnf("failed to create a new vulkan physical device")
		}
		logger.DefaultLogger.Debug("created new vulkan physical device")
		if physical.IsVulkanPhysicalDeviceSupported(&device) {
			devices.PhysicalDevices = append(devices.PhysicalDevices, device)
		}
	}
	slices.SortFunc(devices.PhysicalDevices, func(a, b physical.VulkanPhysicalDevice) int {
		return physical.CompareVulkanPhysicalDevice(&b) - physical.CompareVulkanPhysicalDevice(&a)
	})
	logger.DefaultLogger.Debugf("found %d/%d suitable physical devices", len(devices.PhysicalDevices), len(rawDevices))
	if len(devices.PhysicalDevices) == 0 {
		return devices, errors.New("no suitable physical device found")
	}

	// this will work, as long as the array doesn't relocate :)
	if devices.LogicalDevice, err = logical.NewVulkanLogicalDevice(&devices.PhysicalDevices[0]); err != nil {
		return devices, err
	}

	return devices, nil
}

func FreeVulkanInstanceDevices(devices *VulkanInstanceDevices) error {
	return logical.FreeVulkanLogicalDevice(&devices.LogicalDevice)
}
