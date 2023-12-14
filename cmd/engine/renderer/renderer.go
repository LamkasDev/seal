package renderer

import (
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/vulkan-go/vulkan"
)

type Renderer struct {
	Options   RendererOptions
	Surface   vulkan.Surface
	Swapchain sealVulkan.VulkanSwapchain
}

func NewRenderer(options RendererOptions) (Renderer, error) {
	var err error
	renderer := Renderer{
		Options: options,
	}

	var surfaceRaw uintptr
	if surfaceRaw, err = renderer.Options.Window.Handle.CreateWindowSurface(renderer.Options.VulkanInstance.Handle, nil); err != nil {
		return renderer, err
	}
	renderer.Surface = vulkan.Surface(vulkan.SurfaceFromPointer(surfaceRaw))

	if err := sealVulkan.InitializeVulkanInstanceDevices(renderer.Options.VulkanInstance, options.Window.Handle, &renderer.Surface); err != nil {
		return renderer, err
	}

	if renderer.Swapchain, err = sealVulkan.NewVulkanSwapchain(&renderer.Options.VulkanInstance.Devices.LogicalDevice, &renderer.Surface); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func RunRenderer(renderer *Renderer) error {
	return nil
}

func FreeRenderer(renderer *Renderer) error {
	vulkan.DestroySurface(renderer.Options.VulkanInstance.Handle, renderer.Surface, nil)
	if err := sealVulkan.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}

	return nil
}
