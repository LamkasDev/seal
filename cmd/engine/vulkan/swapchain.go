package vulkan

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchain struct {
	Handle  vulkan.Swapchain
	Device  *VulkanLogicalDevice
	Options VulkanSwapchainOptions
	Images  []VulkanImage
}

func NewVulkanSwapchain(device *VulkanLogicalDevice, surface *vulkan.Surface) (VulkanSwapchain, error) {
	var err error
	swapchain := VulkanSwapchain{
		Device: device,
		Images: []VulkanImage{},
	}

	if swapchain.Options, err = NewVulkanSwapchainOptions(device, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan swapchain options")

	var vulkanSwapchain vulkan.Swapchain
	if res := vulkan.CreateSwapchain(device.Handle, &swapchain.Options.CreateInfo, nil, &vulkanSwapchain); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	swapchain.Handle = vulkanSwapchain
	logger.DefaultLogger.Debug("created new vulkan swapchain")

	var rawImagesCount uint32
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &rawImagesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	rawImages := make([]vulkan.Image, rawImagesCount)
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &rawImagesCount, rawImages); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	for i := 0; i < len(rawImages); i++ {
		image, err := NewVulkanImage(device, rawImages[i])
		if err != nil {
			return swapchain, err
		}
		swapchain.Images = append(swapchain.Images, image)
	}

	return swapchain, nil
}

func FreeVulkanSwapchain(swapchain *VulkanSwapchain) error {
	vulkan.DestroySwapchain(swapchain.Device.Handle, swapchain.Handle, nil)
	for _, image := range swapchain.Images {
		if err := FreeVulkanImage(&image); err != nil {
			return err
		}
	}
	return nil
}
