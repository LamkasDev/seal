package engine

import (
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"

	"github.com/LamkasDev/seal/cmd/common/arch"
	sealGLFW "github.com/LamkasDev/seal/cmd/engine/glfw"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	sealVulkan "github.com/LamkasDev/seal/cmd/engine/vulkan"
	"github.com/LamkasDev/seal/cmd/logger"
)

var CpuProfileFile *os.File
var CpuImageFile *os.File
var MemProfileFile *os.File
var MemImageFile *os.File

func InitializeEngine() {
	var err error
	if arch.SealDebug {
		runtime.MemProfileRate = 1
	}

	// Initialize libraries
	if err = logger.StartLogger(); err != nil {
		panic(err)
	}

	progress.AdvanceLoading()
	if err = sealGLFW.StartGLFW(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	progress.AdvanceLoading()
	if err = sealVulkan.StartVulkan(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	// Start profiling
	if arch.SealDebug {
		CpuProfileFile, err = os.Create("../cpu.prof")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		CpuImageFile, err = os.Create("../cpu.svg")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		if err = pprof.StartCPUProfile(CpuProfileFile); err != nil {
			logger.DefaultLogger.DPanic(err)
		}

		MemProfileFile, err = os.Create("../mem.prof")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
		MemImageFile, err = os.Create("../mem.svg")
		if err != nil {
			logger.DefaultLogger.DPanic(err)
		}
	}
}

func ShutdownEngine() {
	pprof.StopCPUProfile()
	cmd := exec.Command("go", "tool", "pprof", "-svg", "seal_engine.exe", "../cpu.prof")
	cmd.Stdout = CpuImageFile
	if err := cmd.Run(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}
	if err := CpuProfileFile.Close(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}
	if err := CpuImageFile.Close(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}

	if err := pprof.WriteHeapProfile(MemProfileFile); err != nil {
		logger.DefaultLogger.DPanic(err)
	}
	cmd = exec.Command("go", "tool", "pprof", "-svg", "seal_engine.exe", "../mem.prof")
	cmd.Stdout = MemImageFile
	if err := cmd.Run(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}
	if err := MemProfileFile.Close(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}
	if err := MemImageFile.Close(); err != nil {
		logger.DefaultLogger.DPanic(err)
	}

	sealVulkan.EndVulkan()
	sealGLFW.EndGLFW()
	logger.EndLogger()
}
