package progress

import "github.com/LamkasDev/seal/cmd/logger"

const (
	LOADING_STAGE_STARTING_GLFW            = "starting glfw"
	LOADING_STAGE_STARTING_VULKAN          = "starting vulkan"
	LOADING_STAGE_CREATING_VULKAN_INSTANCE = "creating vulkan instance"
	LOADING_STAGE_CREATING_WINDOW          = "creating window"
	LOADING_STAGE_CREATING_ENGINE          = "creating engine"
)

var LoadingStages = []string{
	LOADING_STAGE_STARTING_GLFW,
	LOADING_STAGE_STARTING_VULKAN,
	LOADING_STAGE_CREATING_VULKAN_INSTANCE,
	LOADING_STAGE_CREATING_WINDOW,
	LOADING_STAGE_CREATING_ENGINE,
}
var LoadingStagesLength = len(LoadingStages)

var LoadingStage = 0

func AdvanceLoading() {
	logger.DefaultLogger.Infof("[%d/%d] %s...", LoadingStage+1, LoadingStagesLength, LoadingStages[LoadingStage])
	LoadingStage++
}
