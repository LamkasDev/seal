package vulkan

import (
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstance struct {
	Handle       vulkan.Instance
	Capabilities VulkanInstanceCapabilities
	Options      VulkanInstanceOptions
	Devices      VulkanInstanceDevices
}

func NewVulkanInstance() (VulkanInstance, error) {
	var err error
	instance := VulkanInstance{}

	if instance.Capabilities, err = NewVulkanInstanceCapabilities(); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance capabilities")

	if instance.Options, err = NewVulkanInstanceOptions(&instance); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance options")

	var vulkanInstance vulkan.Instance
	vulkan.CreateInstance(&instance.Options.CreateInfo, nil, &vulkanInstance)
	instance.Handle = vulkanInstance
	logger.DefaultLogger.Info("created new vulkan instance")

	return instance, nil
}

func InitializeVulkanInstanceDevices(instance *VulkanInstance, window *glfw.Window, surface *vulkan.Surface) error {
	var err error
	if instance.Devices, err = NewVulkanInstanceDevices(instance, window, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance devices")

	return nil
}

func FreeVulkanInstance(instance *VulkanInstance) error {
	if err := FreeVulkanInstanceDevices(&instance.Devices); err != nil {
		return err
	}
	vulkan.DestroyInstance(instance.Handle, nil)
	return nil
}
