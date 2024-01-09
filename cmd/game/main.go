package main

import (
	"runtime"

	"github.com/EngoEngine/glm"
	sealEngine "github.com/LamkasDev/seal/cmd/engine/engine"
	"github.com/LamkasDev/seal/cmd/engine/entity"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	"github.com/LamkasDev/seal/cmd/engine/scene"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/font"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/renderer"
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
	if _, err = scene.SpawnSceneModel(&engine.Scene, engine.Renderer.MeshContainer.Meshes[mesh.MESH_BASIC], sealTransform.VulkanTransform{Position: glm.Vec3{0, 0, 0}}); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	uif, _ := font.CreateVulkanFontMesh(&engine.Renderer.MeshContainer, engine.Renderer.FontContainer.Fonts[font.FONT_DEFAULT])
	uie, _ := scene.SpawnSceneModel(&engine.Scene, &uif, sealTransform.VulkanTransform{Position: glm.Vec3{0, 0, 0}, Rotation: glm.Vec3{0, 0, 90}})
	uie.Layer = entity.LAYER_UI

	// Run the main loop
	if err := renderer.PushVulkanRendererBuffers(&engine.Renderer); err != nil {
		return
	}
	if err := sealEngine.RunEngine(&engine); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
}
