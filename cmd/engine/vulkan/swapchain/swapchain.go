package swapchain

import (
	sealFramebuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/framebuffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchain struct {
	Handle       vulkan.Swapchain
	Device       *logical.VulkanLogicalDevice
	Pipeline     *pipeline.VulkanPipeline
	Options      VulkanSwapchainOptions
	Images       []vulkan.Image
	Framebuffers []sealFramebuffer.VulkanFramebuffer
}

func NewVulkanSwapchain(device *logical.VulkanLogicalDevice, pipeline *pipeline.VulkanPipeline, surface *vulkan.Surface) (VulkanSwapchain, error) {
	var err error
	swapchain := VulkanSwapchain{
		Device:       device,
		Pipeline:     pipeline,
		Framebuffers: []sealFramebuffer.VulkanFramebuffer{},
	}

	if swapchain.Options, err = NewVulkanSwapchainOptions(device, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	var vulkanSwapchain vulkan.Swapchain
	if res := vulkan.CreateSwapchain(device.Handle, &swapchain.Options.CreateInfo, nil, &vulkanSwapchain); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	swapchain.Handle = vulkanSwapchain
	logger.DefaultLogger.Debug("created new vulkan swapchain")

	var imagesCount uint32
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &imagesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	swapchain.Images = make([]vulkan.Image, imagesCount)
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &imagesCount, swapchain.Images); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}

	for i := 0; i < len(swapchain.Images); i++ {
		imageview, err := image.NewVulkanImageView(device, &swapchain.Images[i])
		if err != nil {
			return swapchain, err
		}
		framebuffer, err := sealFramebuffer.NewVulkanFramebuffer(device, &swapchain.Pipeline.RenderPass, &imageview, swapchain.Options.CreateInfo.ImageExtent)
		if err != nil {
			return swapchain, err
		}
		swapchain.Framebuffers = append(swapchain.Framebuffers, framebuffer)
	}

	return swapchain, nil
}

func FreeVulkanSwapchain(swapchain *VulkanSwapchain) error {
	for _, framebuffer := range swapchain.Framebuffers {
		if err := sealFramebuffer.FreeVulkanFramebuffer(&framebuffer); err != nil {
			return err
		}
	}
	vulkan.DestroySwapchain(swapchain.Device.Handle, swapchain.Handle, nil)

	return nil
}
