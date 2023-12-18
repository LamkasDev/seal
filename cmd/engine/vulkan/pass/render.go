package pass

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanRenderPass struct {
	Handle  vulkan.RenderPass
	Device  *logical.VulkanLogicalDevice
	Options VulkanRenderPassOptions
}

func NewVulkanRenderPass(device *logical.VulkanLogicalDevice, format vulkan.Format) (VulkanRenderPass, error) {
	pass := VulkanRenderPass{
		Device:  device,
		Options: NewVulkanRenderPassOptions(format),
	}

	var vulkanRenderPass vulkan.RenderPass
	if res := vulkan.CreateRenderPass(device.Handle, &pass.Options.CreateInfo, nil, &vulkanRenderPass); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return pass, vulkan.Error(res)
	}
	pass.Handle = vulkanRenderPass
	logger.DefaultLogger.Debug("created new vulkan render pass")

	return pass, nil
}

func FreeVulkanRenderPass(pass *VulkanRenderPass) error {
	vulkan.DestroyRenderPass(pass.Device.Handle, pass.Handle, nil)
	return nil
}
