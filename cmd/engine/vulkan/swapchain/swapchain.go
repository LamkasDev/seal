package swapchain

import (
	sealFramebuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/framebuffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	sealWindow "github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchain struct {
	Handle     vulkan.Swapchain
	Device     *logical.VulkanLogicalDevice
	Window     *sealWindow.Window
	RenderPass *vulkan.RenderPass
	Options    VulkanSwapchainOptions

	Images         []vulkan.Image
	DepthImage     image.VulkanImage
	DepthImageView image.VulkanImageView
	Framebuffers   []sealFramebuffer.VulkanFramebuffer
}

func NewVulkanSwapchain(device *logical.VulkanLogicalDevice, window *sealWindow.Window, renderPass *pass.VulkanRenderPass, surface *vulkan.Surface, old *VulkanSwapchain) (VulkanSwapchain, error) {
	var err error
	swapchain := VulkanSwapchain{
		Device:       device,
		Window:       window,
		Options:      NewVulkanSwapchainOptions(device, window, surface, old),
		Framebuffers: []sealFramebuffer.VulkanFramebuffer{},
	}

	var vulkanSwapchain vulkan.Swapchain
	if res := vulkan.CreateSwapchain(device.Handle, &swapchain.Options.CreateInfo, nil, &vulkanSwapchain); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	swapchain.Handle = vulkanSwapchain
	logger.DefaultLogger.Debug("created new vulkan swapchain")

	var imagesCount uint32
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &imagesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	swapchain.Images = make([]vulkan.Image, imagesCount)
	swapchain.Framebuffers = make([]sealFramebuffer.VulkanFramebuffer, imagesCount)
	if res := vulkan.GetSwapchainImages(device.Handle, swapchain.Handle, &imagesCount, swapchain.Images); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}

	if swapchain.DepthImage, err = image.NewVulkanImage(device, vulkan.FormatD32Sfloat, window.Data.Extent.Width, window.Data.Extent.Height, vulkan.ImageUsageFlags(vulkan.ImageUsageDepthStencilAttachmentBit)); err != nil {
		return swapchain, err
	}
	if swapchain.DepthImageView, err = image.NewVulkanImageView(device, &swapchain.DepthImage, vulkan.ImageAspectFlags(vulkan.ImageAspectDepthBit)); err != nil {
		return swapchain, err
	}

	for i := 0; i < len(swapchain.Images); i++ {
		imageView, err := image.NewVulkanImageViewRaw(swapchain.Device, &swapchain.Images[i], swapchain.Device.Physical.Capabilities.Surface.ImageFormats[swapchain.Device.Physical.Capabilities.Surface.ImageFormatIndex].Format, vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit))
		if err != nil {
			return swapchain, err
		}
		framebuffer, err := sealFramebuffer.NewVulkanFramebuffer(device, renderPass, &imageView, &swapchain.DepthImageView)
		if err != nil {
			return swapchain, err
		}
		swapchain.Framebuffers[i] = framebuffer
	}

	return swapchain, nil
}

func FreeVulkanSwapchain(swapchain *VulkanSwapchain) error {
	vulkan.DeviceWaitIdle(swapchain.Device.Handle)
	for _, framebuffer := range swapchain.Framebuffers {
		if err := sealFramebuffer.FreeVulkanFramebuffer(&framebuffer); err != nil {
			return err
		}
	}
	if err := image.FreeVulkanImageView(&swapchain.DepthImageView); err != nil {
		return err
	}
	if err := image.FreeVulkanImage(&swapchain.DepthImage); err != nil {
		return err
	}
	vulkan.DestroySwapchain(swapchain.Device.Handle, swapchain.Handle, nil)

	return nil
}
