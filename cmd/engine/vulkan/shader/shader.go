package shader

import (
	"fmt"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanShader struct {
	Vertex   VulkanShaderModule
	Fragment VulkanShaderModule
}

func NewVulkanShader(device *logical.VulkanLogicalDevice, id string) (VulkanShader, error) {
	var err error
	shader := VulkanShader{}

	if shader.Vertex, err = NewVulkanShaderModule(device, fmt.Sprintf("%s.vert", id), vulkan.ShaderStageVertexBit); err != nil {
		return shader, err
	}
	if shader.Fragment, err = NewVulkanShaderModule(device, fmt.Sprintf("%s.frag", id), vulkan.ShaderStageFragmentBit); err != nil {
		return shader, err
	}
	logger.DefaultLogger.Debug("created new vulkan shader")

	return shader, nil
}

func FreeVulkanShader(device *logical.VulkanLogicalDevice, shader *VulkanShader) error {
	if err := FreeVulkanShaderModule(device, &shader.Vertex); err != nil {
		return err
	}
	if err := FreeVulkanShaderModule(device, &shader.Fragment); err != nil {
		return err
	}

	return nil
}
