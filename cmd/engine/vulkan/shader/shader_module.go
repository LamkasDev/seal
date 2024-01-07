package shader

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanShaderModule struct {
	Handle  vulkan.ShaderModule
	Stage   vulkan.PipelineShaderStageCreateInfo
	Options VulkanShaderModuleOptions
}

func CompileVulkanShaderModule(source string, target string) error {
	sourceStat, sourceErr := os.Stat(source)
	targetStat, targetErr := os.Stat(target)
	if sourceErr == nil && targetErr == nil && targetStat.ModTime().After(sourceStat.ModTime()) {
		return nil
	}
	logger.DefaultLogger.Debugf("compiling shader module '%s'...", source)
	cmd := exec.Command("../../resources/glslc", source, "-o", target)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func NewVulkanShaderModule(device *logical.VulkanLogicalDevice, id string, stage vulkan.ShaderStageFlagBits) (VulkanShaderModule, error) {
	var err error
	shaderModule := VulkanShaderModule{}

	source := fmt.Sprintf("../../resources/shaders/%s", id)
	target := fmt.Sprintf("../shaders/%s.spv", id)
	if arch.SealDebug {
		if err = CompileVulkanShaderModule(source, target); err != nil {
			return shaderModule, err
		}
	}

	if shaderModule.Options, err = NewVulkanShaderModuleOptions(target); err != nil {
		return shaderModule, err
	}

	var vulkanShaderModule vulkan.ShaderModule
	if res := vulkan.CreateShaderModule(device.Handle, &shaderModule.Options.CreateInfo, nil, &vulkanShaderModule); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	shaderModule.Handle = vulkanShaderModule
	logger.DefaultLogger.Debug("created new vulkan shader module")

	shaderModule.Stage = vulkan.PipelineShaderStageCreateInfo{
		SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
		Stage:  stage,
		Module: shaderModule.Handle,
		PName:  "main\x00",
	}

	return shaderModule, nil
}

func FreeVulkanShaderModule(device *logical.VulkanLogicalDevice, module *VulkanShaderModule) error {
	vulkan.DestroyShaderModule(device.Handle, module.Handle, nil)
	return nil
}
