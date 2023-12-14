package vulkan

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanImage struct {
	Handle    vulkan.Image
	ImageView VulkanImageView
}

func NewVulkanImage(device *VulkanLogicalDevice, handle vulkan.Image) (VulkanImage, error) {
	var err error
	image := VulkanImage{
		Handle: handle,
	}

	if image.ImageView, err = NewVulkanImageView(device, &image.Handle); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	return image, nil
}

func FreeVulkanImage(image *VulkanImage) error {
	return FreeVulkanImageView(&image.ImageView)
}
