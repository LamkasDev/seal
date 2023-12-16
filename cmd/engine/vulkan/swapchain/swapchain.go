package swapchain

import (
	sealImage "github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchain struct {
	Handle   vulkan.Swapchain
	Device   *logical.VulkanLogicalDevice
	Options  VulkanSwapchainOptions
	Images   []sealImage.VulkanImage
	Pipeline pipeline.VulkanPipeline
}

func NewVulkanSwapchain(device *logical.VulkanLogicalDevice, surface *vulkan.Surface, container *shader.VulkanShaderContainer) (VulkanSwapchain, error) {
	var err error
	swapchain := VulkanSwapchain{
		Device: device,
		Images: []sealImage.VulkanImage{},
	}

	if swapchain.Options, err = NewVulkanSwapchainOptions(device, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan swapchain options")

	if swapchain.Pipeline, err = pipeline.NewVulkanPipeline(device, swapchain.Options.CreateInfo.ImageFormat, swapchain.Options.CreateInfo.ImageExtent, container); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan swapchain pipeline")

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
		image, err := sealImage.NewVulkanImage(device, rawImages[i])
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
		if err := sealImage.FreeVulkanImage(&image); err != nil {
			return err
		}
	}
	pipeline.FreeVulkanPipeline(&swapchain.Pipeline)

	return nil
}
