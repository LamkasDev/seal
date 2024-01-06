package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanImage struct {
	Handle       vulkan.Image
	MemoryHandle vulkan.DeviceMemory
	Device       *logical.VulkanLogicalDevice
	Options      VulkanImageOptions
}

func NewVulkanImage(device *logical.VulkanLogicalDevice, w uint32, h uint32) (VulkanImage, error) {
	image := VulkanImage{
		Device:  device,
		Options: NewVulkanImageOptions(w, h),
	}

	var vulkanImage vulkan.Image
	if res := vulkan.CreateImage(device.Handle, &image.Options.CreateInfo, nil, &vulkanImage); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return image, vulkan.Error(res)
	}
	image.Handle = vulkanImage
	logger.DefaultLogger.Debug("created new vulkan image")

	if err := AllocateVulkanImage(&image); err != nil {
		return image, err
	}
	logger.DefaultLogger.Debug("allocated vulkan image memory")

	return image, nil
}

func AllocateVulkanImage(image *VulkanImage) error {
	var requirements vulkan.MemoryRequirements
	vulkan.GetImageMemoryRequirements(image.Device.Handle, image.Handle, &requirements)
	requirements.Deref()

	if err := UpdateVulkanImageOptions(&image.Options, image.Device, requirements); err != nil {
		return err
	}

	var vulkanImageMemory vulkan.DeviceMemory
	if res := vulkan.AllocateMemory(image.Device.Handle, &image.Options.AllocateInfo, nil, &vulkanImageMemory); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	image.MemoryHandle = vulkanImageMemory

	if res := vulkan.BindImageMemory(image.Device.Handle, image.Handle, image.MemoryHandle, 0); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}

	return nil
}

func FreeVulkanImage(image *VulkanImage) error {
	vulkan.FreeMemory(image.Device.Handle, image.MemoryHandle, nil)
	vulkan.DestroyImage(image.Device.Handle, image.Handle, nil)
	return nil
}
