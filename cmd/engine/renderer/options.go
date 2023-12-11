package renderer

import (
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/engine/window"
)

type RendererOptions struct {
	VulkanInstance *sealVulkan.VulkanInstance
	Window         *window.Window
}

func NewRendererOptions(vulkanInstance *sealVulkan.VulkanInstance, window *window.Window) RendererOptions {
	return RendererOptions{
		VulkanInstance: vulkanInstance,
		Window:         window,
	}
}
