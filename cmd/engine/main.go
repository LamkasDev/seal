package main

import (
	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/LamkasDev/seal/cmd/engine/engine"
	sealGLFW "github.com/LamkasDev/seal/cmd/engine/glfw"
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

	if err := sealGLFW.StartGLFW(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Info("started glfw")
	defer sealGLFW.EndGLFW()

	// Setup vulkan loader
	vulkanLoader, err := sealVulkan.NewVulkanLoader()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Info("create new vulkan loader")
	defer sealVulkan.FreeVulkanLoader(vulkanLoader)

	// Setup required instances
	vulkanInstance, err := sealVulkan.NewVulkanInstance(vulkanLoader)
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Info("created new vulkan instance")
	defer sealVulkan.FreeVulkanInstance(&vulkanInstance)

	sealWindow, err := window.NewWindow(window.NewWindowOptions("Test", constants.DefaultResolution))
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Info("created new window")
	defer window.FreeWindow(&sealWindow)

	sealEngine, err := engine.NewEngine(&sealWindow)
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Info("created new engine")
	defer engine.FreeEngine(&sealEngine)

	// Run the main loop
	if err := engine.RunEngine(&sealEngine); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}
