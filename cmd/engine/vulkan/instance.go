package vulkan

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/device"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstance struct {
	Handle       vulkan.Instance
	Capabilities VulkanInstanceCapabilities
	Options      VulkanInstanceOptions
	Devices      device.VulkanInstanceDevices
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
	if res := vulkan.CreateInstance(&instance.Options.CreateInfo, nil, &vulkanInstance); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	instance.Handle = vulkanInstance
	logger.DefaultLogger.Debug("created new vulkan instance")

	return instance, nil
}

func InitializeVulkanInstanceDevices(instance *VulkanInstance, window *glfw.Window, surface *vulkan.Surface) error {
	var err error
	if instance.Devices, err = device.NewVulkanInstanceDevices(instance.Handle, window, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance devices")

	return nil
}

func FreeVulkanInstance(instance *VulkanInstance) error {
	if err := device.FreeVulkanInstanceDevices(&instance.Devices); err != nil {
		return err
	}
	vulkan.DestroyInstance(instance.Handle, nil)
	return nil
}
