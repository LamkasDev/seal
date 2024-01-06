package shader

import (
	"os"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
)

const SHADER_BASIC = "basic"

var DefaultShaders = []string{
	SHADER_BASIC,
}

type VulkanShaderContainer struct {
	Device  *logical.VulkanLogicalDevice
	Shaders map[string]VulkanShader
}

func NewVulkanShaderContainer(device *logical.VulkanLogicalDevice) (VulkanShaderContainer, error) {
	container := VulkanShaderContainer{
		Device:  device,
		Shaders: map[string]VulkanShader{},
	}

	if err := os.MkdirAll("../shaders", 0755); err != nil {
		logger.DefaultLogger.Error(err)
		return container, err
	}
	for _, shader := range DefaultShaders {
		if _, err := CreateVulkanShaderWithContainer(&container, shader); err != nil {
			return container, err
		}
	}
	logger.DefaultLogger.Debug("created new vulkan shader container")

	return container, nil
}

func CreateVulkanShaderWithContainer(container *VulkanShaderContainer, id string) (VulkanShader, error) {
	shader, err := NewVulkanShader(container.Device, id)
	if err != nil {
		return shader, err
	}
	container.Shaders[id] = shader

	return shader, nil
}

func FreeVulkanShaderContainer(container *VulkanShaderContainer) error {
	for _, shader := range container.Shaders {
		if err := FreeVulkanShader(&shader); err != nil {
			return err
		}
	}

	return nil
}
