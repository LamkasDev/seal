package vulkan

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanPhysicalDevice struct {
	Handle                       vulkan.PhysicalDevice
	Properties                   vulkan.PhysicalDeviceProperties
	Features                     vulkan.PhysicalDeviceFeatures
	QueueFamilies                []vulkan.QueueFamilyProperties
	QueueFamilyGraphicsIndex     int
	QueueFamilyPresentationIndex int
}

func NewVulkanPhysicalDevice(handle vulkan.PhysicalDevice) (VulkanPhysicalDevice, error) {
	device := VulkanPhysicalDevice{
		Handle: handle,
	}
	vulkan.GetPhysicalDeviceProperties(device.Handle, &device.Properties)
	device.Properties.Deref()
	vulkan.GetPhysicalDeviceFeatures(device.Handle, &device.Features)
	device.Features.Deref()

	var queueFamilyPropertiesCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device.Handle, &queueFamilyPropertiesCount, nil)
	device.QueueFamilies = make([]vulkan.QueueFamilyProperties, queueFamilyPropertiesCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(device.Handle, &queueFamilyPropertiesCount, device.QueueFamilies)

	for i := 0; i < len(device.QueueFamilies); i++ {
		device.QueueFamilies[i].Deref()
		if vulkan.QueueFlagBits(device.QueueFamilies[i].QueueFlags)&vulkan.QueueGraphicsBit == vulkan.QueueGraphicsBit {
			device.QueueFamilyGraphicsIndex = i
		}
	}

	return device, nil
}

func ProbeVulkanPhysicalDevice(device *VulkanPhysicalDevice, surface *vulkan.Surface) {
	for i := uint32(0); i < uint32(len(device.QueueFamilies)); i++ {
		var supported vulkan.Bool32
		vulkan.GetPhysicalDeviceSurfaceSupport(device.Handle, i, *surface, &supported)
		if supported == 1 {
			device.QueueFamilyPresentationIndex = int(i)
			break
		}
	}
}

func CompareVulkanPhysicalDevice(device *VulkanPhysicalDevice) int {
	return int(device.Properties.Limits.MaxImageDimension2D)
}

func IsVulkanPhysicalDeviceSupported(device *VulkanPhysicalDevice) bool {
	return (device.Features.GeometryShader == 1) && device.QueueFamilyGraphicsIndex != -1
}
