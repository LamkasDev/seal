package device

import (
	"errors"
	"slices"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceDevices struct {
	PhysicalDevices []physical.VulkanPhysicalDevice
	LogicalDevice   logical.VulkanLogicalDevice
}

func NewVulkanInstanceDevices(instance vulkan.Instance, cwindow *window.Window, surface *vulkan.Surface) (VulkanInstanceDevices, error) {
	var err error
	devices := VulkanInstanceDevices{}

	var rawDevicesCount uint32
	if res := vulkan.EnumeratePhysicalDevices(instance, &rawDevicesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return devices, vulkan.Error(res)
	}
	rawDevices := make([]vulkan.PhysicalDevice, rawDevicesCount)
	if res := vulkan.EnumeratePhysicalDevices(instance, &rawDevicesCount, rawDevices); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return devices, vulkan.Error(res)
	}

	devices.PhysicalDevices = []physical.VulkanPhysicalDevice{}
	for i := 0; i < len(rawDevices); i++ {
		var device physical.VulkanPhysicalDevice
		if device, err = physical.NewVulkanPhysicalDevice(rawDevices[i], cwindow, surface); err != nil {
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

	if devices.LogicalDevice, err = logical.NewVulkanLogicalDevice(&devices.PhysicalDevices[0]); err != nil {
		return devices, err
	}

	return devices, nil
}

func UpdateVulkanInstanceDevices(devices *VulkanInstanceDevices) error {
	for i := 0; i < len(devices.PhysicalDevices); i++ {
		if err := physical.UpdateVulkanPhysicalDevice(&devices.PhysicalDevices[i]); err != nil {
			return err
		}
	}

	return nil
}

func FreeVulkanInstanceDevices(devices *VulkanInstanceDevices) error {
	return logical.FreeVulkanLogicalDevice(&devices.LogicalDevice)
}
