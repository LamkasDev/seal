package main

import (
	"runtime"

	"github.com/EngoEngine/glm"
	sealEngine "github.com/LamkasDev/seal/cmd/engine/engine"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/scene"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	sealTransform "github.com/LamkasDev/seal/cmd/engine/vulkan/transform"
	"github.com/LamkasDev/seal/cmd/logger"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	// Initialize engine dependencies
	sealEngine.InitializeEngine()
	defer sealEngine.ShutdownEngine()

	// Setup engine
	progress.AdvanceLoading()
	engine, err := sealEngine.NewEngine()
	if err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	renderer.RendererInstance = &engine.Renderer
	defer sealEngine.FreeEngine(&engine)

	// Do goofy stuff
	if err = scene.SpawnSceneModel(&engine.Scene, engine.Renderer.Pipeline.MeshContainer.Meshes[mesh.MESH_BASIC], sealTransform.VulkanTransform{Position: glm.Vec3{0, 0, 0}}); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}

	// Run the main loop
	if err := sealEngine.RunEngine(&engine); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}
