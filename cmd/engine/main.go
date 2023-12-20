package main

import (
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"

	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/engine/engine"
	sealGLFW "github.com/LamkasDev/seal/cmd/engine/glfw"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/logger"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if arch.SealDebug {
		runtime.MemProfileRate = 1
	}

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
		cpuProfileFile, err := os.Create("../cpu.prof")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		cpuImageFile, err := os.Create("../cpu.svg")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		if err = pprof.StartCPUProfile(cpuProfileFile); err != nil {
			logger.DefaultLogger.DPanic(err)
		}

		memProfileFile, err := os.Create("../mem.prof")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		memImageFile, err := os.Create("../mem.svg")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}

		defer func() {
			pprof.StopCPUProfile()
			cmd := exec.Command("go", "tool", "pprof", "-svg", "seal_engine.exe", "../cpu.prof")
			cmd.Stdout = cpuImageFile
			if err := cmd.Run(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
			if err := cpuProfileFile.Close(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
			if err := cpuImageFile.Close(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}

			if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
			cmd = exec.Command("go", "tool", "pprof", "-svg", "seal_engine.exe", "../mem.prof")
			cmd.Stdout = memImageFile
			if err := cmd.Run(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
			if err := memProfileFile.Close(); err != nil {
				logger.DefaultLogger.DPanic(err)
			}
			if err := memImageFile.Close(); err != nil {
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
