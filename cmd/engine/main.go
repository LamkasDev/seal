package main

import (
	"github.com/LamkasDev/seal/cmd/engine/engine"
	sealGLFW "github.com/LamkasDev/seal/cmd/engine/glfw"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
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

	// Setup engine
	progress.AdvanceLoading()
	sealEngine, err := engine.NewEngine()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	defer engine.FreeEngine(&sealEngine)

	// Run the main loop
	if err := engine.RunEngine(&sealEngine); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}
