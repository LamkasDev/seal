package renderer

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/fence"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/semaphore"
	"github.com/vulkan-go/vulkan"
)

type VulkanRendererSyncer struct {
	Device                   *logical.VulkanLogicalDevice
	ImageAvailableSemaphores []semaphore.VulkanSemaphore
	RenderFinishedSemaphores []semaphore.VulkanSemaphore
	InFlightFences           []fence.VulkanFence
}

func NewVulkanRendererSyncer(device *logical.VulkanLogicalDevice) (VulkanRendererSyncer, error) {
	var err error
	syncer := VulkanRendererSyncer{
		Device:                   device,
		ImageAvailableSemaphores: make([]semaphore.VulkanSemaphore, commonPipeline.MaxFramesInFlight),
		RenderFinishedSemaphores: make([]semaphore.VulkanSemaphore, commonPipeline.MaxFramesInFlight),
		InFlightFences:           make([]fence.VulkanFence, commonPipeline.MaxFramesInFlight),
	}

	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
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

func FreeVulkanRendererSyncer(syncer *VulkanRendererSyncer) error {
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
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
