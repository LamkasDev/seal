package progress

import "github.com/LamkasDev/seal/cmd/logger"

const (
	LOADING_STAGE_STARTING_GLFW                      = "starting glfw"
	LOADING_STAGE_STARTING_VULKAN                    = "starting vulkan"
	LOADING_STAGE_CREATING_ENGINE                    = "creating engine"
	LOADING_STAGE_CREATING_RENDERER                  = "creating renderer"
	LOADING_STAGE_CREATING_VULKAN_INSTANCE           = "renderer > creating vulkan instance"
	LOADING_STAGE_CREATING_WINDOW                    = "renderer > creating window"
	LOADING_STAGE_CREATING_SURFACE                   = "renderer > creating surface"
	LOADING_STAGE_CREATING_DEVICES                   = "renderer > creating devices"
	LOADING_STAGE_CREATING_TEXTURE_CONTAINER         = "renderer > creating texture container"
	LOADING_STAGE_CREATING_SHADER_CONTAINER          = "renderer > creating shader container"
	LOADING_STAGE_CREATING_MESH_CONTAINER            = "renderer > creating mesh container"
	LOADING_STAGE_CREATING_DESCRIPTOR_POOL_CONTAINER = "renderer > creating descriptor pool container"
	LOADING_STAGE_CREATING_BUFFER_CONTAINER          = "renderer > creating buffer container"
	LOADING_STAGE_CREATING_SAMPLER                   = "renderer > creating sampler"
	LOADING_STAGE_CREATING_RENDER_PASS               = "renderer > creating render pass"
	LOADING_STAGE_CREATING_RENDERER_SYNCER           = "renderer > creating renderer syncer"
	LOADING_STAGE_CREATING_RENDERER_COMMANDER        = "renderer > creating renderer commander"
	LOADING_STAGE_CREATING_PIPELINE                  = "renderer > creating pipeline"
	LOADING_STAGE_CREATING_SWAPCHAIN                 = "renderer > creating swapchain"
	LOADING_STAGE_CREATING_INPUT                     = "creating input"
	LOADING_STAGE_CREATING_SCENE                     = "creating scene"
)

var LoadingStages = []string{
	LOADING_STAGE_STARTING_GLFW,
	LOADING_STAGE_STARTING_VULKAN,
	LOADING_STAGE_CREATING_ENGINE,
	LOADING_STAGE_CREATING_RENDERER,
	LOADING_STAGE_CREATING_VULKAN_INSTANCE,
	LOADING_STAGE_CREATING_WINDOW,
	LOADING_STAGE_CREATING_SURFACE,
	LOADING_STAGE_CREATING_DEVICES,
	LOADING_STAGE_CREATING_TEXTURE_CONTAINER,
	LOADING_STAGE_CREATING_SHADER_CONTAINER,
	LOADING_STAGE_CREATING_MESH_CONTAINER,
	LOADING_STAGE_CREATING_DESCRIPTOR_POOL_CONTAINER,
	LOADING_STAGE_CREATING_BUFFER_CONTAINER,
	LOADING_STAGE_CREATING_SAMPLER,
	LOADING_STAGE_CREATING_RENDER_PASS,
	LOADING_STAGE_CREATING_RENDERER_SYNCER,
	LOADING_STAGE_CREATING_RENDERER_COMMANDER,
	LOADING_STAGE_CREATING_PIPELINE,
	LOADING_STAGE_CREATING_SWAPCHAIN,
	LOADING_STAGE_CREATING_SCENE,
	LOADING_STAGE_CREATING_INPUT,
}
var LoadingStagesLength = len(LoadingStages)

var LoadingStage = 0

func AdvanceLoading() {
	logger.DefaultLogger.Infof("[%d/%d] %s...", LoadingStage+1, LoadingStagesLength, LoadingStages[LoadingStage])
	LoadingStage++
}
