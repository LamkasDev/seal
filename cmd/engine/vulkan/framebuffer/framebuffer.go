package framebuffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanFramebuffer struct {
	Handle    vulkan.Framebuffer
	Device    *logical.VulkanLogicalDevice
	Imageview *image.VulkanImageView
	Options   VulkanFramebufferOptions
}

func NewVulkanFramebuffer(device *logical.VulkanLogicalDevice, pass *pass.VulkanRenderPass, imageview *image.VulkanImageView, extent vulkan.Extent2D) (VulkanFramebuffer, error) {
	var err error
	framebuffer := VulkanFramebuffer{
		Device:    device,
		Imageview: imageview,
	}

	if framebuffer.Options, err = NewVulkanFramebufferOptions(pass, imageview, extent); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	var vulkanFramebuffer vulkan.Framebuffer
	if res := vulkan.CreateFramebuffer(device.Handle, &framebuffer.Options.CreateInfo, nil, &vulkanFramebuffer); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	framebuffer.Handle = vulkanFramebuffer
	logger.DefaultLogger.Debug("created new vulkan framebuffer")

	return framebuffer, nil
}

func FreeVulkanFramebuffer(framebuffer *VulkanFramebuffer) error {
	if err := image.FreeVulkanImageView(framebuffer.Imageview); err != nil {
		return err
	}
	vulkan.DestroyFramebuffer(framebuffer.Device.Handle, framebuffer.Handle, nil)
	return nil
}
