package pass

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanRenderPassContainer struct {
	Device *logical.VulkanLogicalDevice
	Passes []*VulkanRenderPass
}

func NewVulkanRenderPassContainer(device *logical.VulkanLogicalDevice, format vulkan.Format, shaders []string) (VulkanRenderPassContainer, error) {
	container := VulkanRenderPassContainer{
		Device: device,
		Passes: []*VulkanRenderPass{},
	}
	if _, err := CreateVulkanRenderPassWithContainer(&container, NewVulkanSceneRenderPassOptions(format, shaders)); err != nil {
		return container, err
	}
	if _, err := CreateVulkanRenderPassWithContainer(&container, NewVulkanUIRenderPassOptions(format, shaders)); err != nil {
		return container, err
	}
	logger.DefaultLogger.Debug("created new vulkan mesh container")

	return container, nil
}

func CreateVulkanRenderPassWithContainer(container *VulkanRenderPassContainer, options VulkanRenderPassOptions) (VulkanRenderPass, error) {
	renderPass, err := NewVulkanRenderPass(container.Device, options)
	if err != nil {
		return renderPass, err
	}
	container.Passes = append(container.Passes, &renderPass)

	return renderPass, nil
}

func FreeVulkanRenderPassContainer(container *VulkanRenderPassContainer) error {
	for _, mesh := range container.Passes {
		if err := FreeVulkanRenderPass(mesh); err != nil {
			return err
		}
	}

	return nil
}
