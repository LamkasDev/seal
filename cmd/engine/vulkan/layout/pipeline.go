package layout

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineLayout struct {
	Handle  vulkan.PipelineLayout
	Device  *logical.VulkanLogicalDevice
	Options VulkanPipelineLayoutOptions
}

func NewVulkanPipelineLayout(device *logical.VulkanLogicalDevice) (VulkanPipelineLayout, error) {
	var err error
	pipelineLayout := VulkanPipelineLayout{
		Device: device,
	}

	if pipelineLayout.Options, err = NewVulkanPipelineLayoutOptions(); err != nil {
		return pipelineLayout, err
	}

	var vulkanPipelineLayout vulkan.PipelineLayout
	if res := vulkan.CreatePipelineLayout(device.Handle, &pipelineLayout.Options.CreateInfo, nil, &vulkanPipelineLayout); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	pipelineLayout.Handle = vulkanPipelineLayout
	logger.DefaultLogger.Debug("created new vulkan pipeline layout")

	return pipelineLayout, nil
}

func FreeVulkanPipelineLayout(layout *VulkanPipelineLayout) error {
	vulkan.DestroyPipelineLayout(layout.Device.Handle, layout.Handle, nil)
	return nil
}
