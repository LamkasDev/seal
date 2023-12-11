package vulkan

import (
	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vkngwrapper/core/v2"
	vulkan "github.com/vkngwrapper/core/v2/core1_0"
)

type VulkanInstance struct {
	Capabilities VulkanInstanceCapabilities
	Options      VulkanInstanceOptions
	Handle       vulkan.Instance
	Debugger     VulkanDebugger
}

func NewVulkanInstance(loader *core.VulkanLoader) (VulkanInstance, error) {
	var err error
	instance := VulkanInstance{}

	if instance.Capabilities, err = NewVulkanInstanceCapabilities(loader); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance capabilities")

	if instance.Options, err = NewVulkanInstanceOptions(&instance.Capabilities); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance options")

	instance.Handle, _, _ = loader.CreateInstance(nil, instance.Options.CreateInfo)

	if arch.SealDebug {
		if instance.Debugger, err = NewVulkanDebugger(); err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
	}

	return instance, nil
}

func FreeVulkanInstance(instance *VulkanInstance) error {
	instance.Handle.Destroy(nil)

	return nil
}
