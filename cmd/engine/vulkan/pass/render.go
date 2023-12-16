package pass

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanRenderPass struct {
	Handle  vulkan.RenderPass
	Options VulkanRenderPassOptions
}

func NewVulkanRenderPass(device *logical.VulkanLogicalDevice, format vulkan.Format) (VulkanRenderPass, error) {
	var err error
	pass := VulkanRenderPass{}

	if pass.Options, err = NewVulkanRenderPassOptions(format); err != nil {
		return pass, err
	}

	var vulkanRenderPass vulkan.RenderPass
	if res := vulkan.CreateRenderPass(device.Handle, &pass.Options.CreateInfo, nil, &vulkanRenderPass); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	pass.Handle = vulkanRenderPass
	logger.DefaultLogger.Debug("created new vulkan render pass")

	return pass, nil
}
