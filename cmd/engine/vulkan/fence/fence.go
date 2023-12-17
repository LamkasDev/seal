package fence

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanFence struct {
	Handle  vulkan.Fence
	Device  *logical.VulkanLogicalDevice
	Options VulkanFenceOptions
}

func NewVulkanFence(device *logical.VulkanLogicalDevice, flags vulkan.FenceCreateFlags) (VulkanFence, error) {
	var err error
	fence := VulkanFence{
		Device: device,
	}

	if fence.Options, err = NewVulkanFenceOptions(flags); err != nil {
		return fence, err
	}

	var vulkanFence vulkan.Fence
	if res := vulkan.CreateFence(device.Handle, &fence.Options.CreateInfo, nil, &vulkanFence); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	fence.Handle = vulkanFence
	logger.DefaultLogger.Debug("created new vulkan fence")

	return fence, nil
}

func FreeVulkanFence(fence *VulkanFence) error {
	vulkan.DestroyFence(fence.Device.Handle, fence.Handle, nil)
	return nil
}
