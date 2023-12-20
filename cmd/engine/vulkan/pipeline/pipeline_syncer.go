package pipeline

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/fence"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/semaphore"
	"github.com/vulkan-go/vulkan"
)

type VulkanPipelineSyncer struct {
	Device                   *logical.VulkanLogicalDevice
	ImageAvailableSemaphores []semaphore.VulkanSemaphore
	RenderFinishedSemaphores []semaphore.VulkanSemaphore
	InFlightFences           []fence.VulkanFence
}

func NewVulkanPipelineSyncer(device *logical.VulkanLogicalDevice) (VulkanPipelineSyncer, error) {
	var err error
	syncer := VulkanPipelineSyncer{
		Device:                   device,
		ImageAvailableSemaphores: make([]semaphore.VulkanSemaphore, MaxFramesInFlight),
		RenderFinishedSemaphores: make([]semaphore.VulkanSemaphore, MaxFramesInFlight),
		InFlightFences:           make([]fence.VulkanFence, MaxFramesInFlight),
	}

	for i := 0; i < MaxFramesInFlight; i++ {
		if syncer.ImageAvailableSemaphores[i], err = semaphore.NewVulkanSemaphore(device); err != nil {
			return syncer, err
		}
		if syncer.RenderFinishedSemaphores[i], err = semaphore.NewVulkanSemaphore(device); err != nil {
			return syncer, err
		}
		if syncer.InFlightFences[i], err = fence.NewVulkanFence(device, vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit)); err != nil {
			return syncer, err
		}
	}

	return syncer, nil
}

func FreeVulkanPipelineSyncer(syncer *VulkanPipelineSyncer) error {
	for i := 0; i < MaxFramesInFlight; i++ {
		if err := fence.FreeVulkanFence(&syncer.InFlightFences[i]); err != nil {
			return err
		}
		if err := semaphore.FreeVulkanSemaphore(&syncer.RenderFinishedSemaphores[i]); err != nil {
			return err
		}
		if err := semaphore.FreeVulkanSemaphore(&syncer.ImageAvailableSemaphores[i]); err != nil {
			return err
		}
	}

	return nil
}
