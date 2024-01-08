package framebuffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	sealPass "github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanFramebuffer struct {
	Handle    vulkan.Framebuffer
	Device    *logical.VulkanLogicalDevice
	Pass      *sealPass.VulkanRenderPass
	Imageview *image.VulkanImageView
	Options   VulkanFramebufferOptions
}

func NewVulkanFramebuffer(device *logical.VulkanLogicalDevice, pass *sealPass.VulkanRenderPass, imageView *image.VulkanImageView, depthImageView *image.VulkanImageView) (VulkanFramebuffer, error) {
	framebuffer := VulkanFramebuffer{
		Device:    device,
		Pass:      pass,
		Imageview: imageView,
		Options:   NewVulkanFramebufferOptions(device, pass, imageView, depthImageView),
	}

	var vulkanFramebuffer vulkan.Framebuffer
	if res := vulkan.CreateFramebuffer(device.Handle, &framebuffer.Options.CreateInfo, nil, &vulkanFramebuffer); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return framebuffer, vulkan.Error(res)
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
