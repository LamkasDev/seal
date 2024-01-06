package transform_buffer

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	sealUniform "github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/samber/lo"
	"github.com/vulkan-go/vulkan"
)

type VulkanTransformBufferOptions struct {
	Uniforms []sealUniform.VulkanUniform
}

func NewVulkanTransformBufferOptions(uniform sealUniform.VulkanUniform) VulkanTransformBufferOptions {
	data := VulkanTransformBufferOptions{
		Uniforms: lo.RepeatBy(commonPipeline.MaxFramesInFlight, func(index int) sealUniform.VulkanUniform {
			return uniform
		}),
	}

	return data
}

func GetVulkanTransformBufferOptionsUniformsOffset(options *VulkanTransformBufferOptions) vulkan.DeviceSize {
	return 0
}

func GetVulkanTransformBufferOptionsUniformsSize(options *VulkanTransformBufferOptions) vulkan.DeviceSize {
	return vulkan.DeviceSize(int(sealUniform.VulkanUniformSize) * len(options.Uniforms))
}

func GetVulkanTransformBufferOptionsSize(options *VulkanTransformBufferOptions) vulkan.DeviceSize {
	return GetVulkanTransformBufferOptionsUniformsOffset(options) + GetVulkanTransformBufferOptionsUniformsSize(options)
}
