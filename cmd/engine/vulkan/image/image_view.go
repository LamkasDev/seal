package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanImageView struct {
	Handle  vulkan.ImageView
	Device  *logical.VulkanLogicalDevice
	Options VulkanImageViewOptions
}

func NewVulkanImageView(device *logical.VulkanLogicalDevice, image *vulkan.Image) (VulkanImageView, error) {
	imageView := VulkanImageView{
		Device:  device,
		Options: NewVulkanImageViewOptions(device, image),
	}

	var vulkanImageView vulkan.ImageView
	if res := vulkan.CreateImageView(device.Handle, &imageView.Options.CreateInfo, nil, &vulkanImageView); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return imageView, vulkan.Error(res)
	}
	imageView.Handle = vulkanImageView
	logger.DefaultLogger.Debug("created new vulkan image view")

	return imageView, nil
}

func FreeVulkanImageView(imageView *VulkanImageView) error {
	vulkan.DestroyImageView(imageView.Device.Handle, imageView.Handle, nil)
	return nil
}
