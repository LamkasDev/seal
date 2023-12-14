package main

import (
	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/LamkasDev/seal/cmd/engine/engine"
	sealGLFW "github.com/LamkasDev/seal/cmd/engine/glfw"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
)

func main() {
	// Initialize libraries
	if err := logger.StartLogger(); err != nil {
		panic(err)
	}
	defer logger.EndLogger()

	progress.AdvanceLoading()
	if err := sealGLFW.StartGLFW(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer sealGLFW.EndGLFW()

	progress.AdvanceLoading()
	if err := sealVulkan.StartVulkan(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer sealVulkan.EndVulkan()

	// Setup required instances
	progress.AdvanceLoading()
	vulkanInstance, err := sealVulkan.NewVulkanInstance()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer sealVulkan.FreeVulkanInstance(&vulkanInstance)

	progress.AdvanceLoading()
	sealWindow, err := window.NewWindow(window.NewWindowOptions("Test", constants.DefaultResolution))
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer window.FreeWindow(&sealWindow)

	progress.AdvanceLoading()
	sealEngine, err := engine.NewEngine(renderer.NewRendererOptions(&vulkanInstance, &sealWindow))
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer engine.FreeEngine(&sealEngine)

	// Run the main loop
	if err := engine.RunEngine(&sealEngine); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}
