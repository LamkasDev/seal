package main

import (
	"os"
	"os/exec"
	"runtime/pprof"

	"github.com/LamkasDev/seal/cmd/common/arch"
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

	// Start profiling
	if arch.SealDebug {
		profileFile, err := os.Create("../cpu.prof")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		imageFile, err := os.Create("../cpu.svg")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		pprof.StartCPUProfile(profileFile)
		defer func() {
			pprof.StopCPUProfile()
			cmd := exec.Command("go", "tool", "pprof", "-svg", "seal_engine.exe", "../cpu.prof")
			cmd.Stdout = imageFile
			if err := cmd.Run(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
		}()
	}

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
