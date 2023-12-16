package renderer

import (
	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/swapchain"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type Renderer struct {
	VulkanInstance  sealVulkan.VulkanInstance
	Window          window.Window
	Surface         vulkan.Surface
	ShaderContainer shader.VulkanShaderContainer
	Swapchain       swapchain.VulkanSwapchain
}

func NewRenderer() (Renderer, error) {
	var err error
	renderer := Renderer{}

	progress.AdvanceLoading()
	renderer.VulkanInstance, err = sealVulkan.NewVulkanInstance()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	renderer.Window, err = window.NewWindow(window.NewWindowOptions("Test", constants.DefaultResolution))
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	var surfaceRaw uintptr
	if surfaceRaw, err = renderer.Window.Handle.CreateWindowSurface(renderer.VulkanInstance.Handle, nil); err != nil {
		return renderer, err
	}
	renderer.Surface = vulkan.Surface(vulkan.SurfaceFromPointer(surfaceRaw))

	progress.AdvanceLoading()
	if err := sealVulkan.InitializeVulkanInstanceDevices(&renderer.VulkanInstance, renderer.Window.Handle, &renderer.Surface); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.ShaderContainer, err = shader.NewVulkanShaderContainer(&renderer.VulkanInstance.Devices.LogicalDevice); err != nil {
		return renderer, err
	}

	progress.AdvanceLoading()
	if renderer.Swapchain, err = swapchain.NewVulkanSwapchain(&renderer.VulkanInstance.Devices.LogicalDevice, &renderer.Surface, &renderer.ShaderContainer); err != nil {
		return renderer, err
	}

	return renderer, nil
}

func RunRenderer(renderer *Renderer) error {
	return nil
}

func FreeRenderer(renderer *Renderer) error {
	if err := swapchain.FreeVulkanSwapchain(&renderer.Swapchain); err != nil {
		return err
	}
	vulkan.DestroySurface(renderer.VulkanInstance.Handle, renderer.Surface, nil)
	window.FreeWindow(&renderer.Window)
	sealVulkan.FreeVulkanInstance(&renderer.VulkanInstance)

	return nil
}
