package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanImage struct {
	Handle  vulkan.Image
	Device  *logical.VulkanLogicalDevice
	Options VulkanImageOptions
}

func NewVulkanImage(device *logical.VulkanLogicalDevice, image *vulkan.Image) (VulkanImage, error) {
	image := VulkanImage{
		Device:  device,
		Options: NewVulkanImageOptions(device, image),
	}

	var vulkanImage vulkan.Image
	if res := vulkan.CreateImage(device.Handle, &image.Options.CreateInfo, nil, &vulkanImage); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return image, vulkan.Error(res)
	}
	image.Handle = vulkanImage
	logger.DefaultLogger.Debug("created new vulkan image")

	return image, nil
}

func FreeVulkanImage(image *VulkanImage) error {
	vulkan.DestroyImage(image.Device.Handle, image.Handle, nil)
	return nil
}
