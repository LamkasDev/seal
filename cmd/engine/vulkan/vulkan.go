package vulkan

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vkngwrapper/core/v2"
)

func NewVulkanLoader() (*core.VulkanLoader, error) {
	return core.CreateLoaderFromProcAddr(glfw.GetVulkanGetInstanceProcAddress())
}

func FreeVulkanLoader(loader *core.VulkanLoader) error {
	loader.Driver().Destroy()
	return nil
}
