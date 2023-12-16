package shader

import (
	"os"
	"unsafe"

	"github.com/vulkan-go/vulkan"
)

type VulkanShaderModuleOptions struct {
	Contents   []byte
	CreateInfo vulkan.ShaderModuleCreateInfo
}

func NewVulkanShaderModuleOptions(path string) (VulkanShaderModuleOptions, error) {
	var err error
	options := VulkanShaderModuleOptions{}

	if options.Contents, err = os.ReadFile(path); err != nil {
		return options, err
	}

	options.CreateInfo = vulkan.ShaderModuleCreateInfo{
		SType:    vulkan.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(options.Contents)),
		PCode:    (*(*[]uint32)(unsafe.Pointer(&options.Contents))),
	}

	return options, nil
}
