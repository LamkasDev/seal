package vulkan

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

func StartVulkan() error {
	vulkan.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	return vulkan.Init()
}

func EndVulkan() error {
	return nil
}
