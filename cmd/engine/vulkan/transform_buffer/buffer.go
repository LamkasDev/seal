package transform_buffer

import (
	"unsafe"

	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	sealPipeline "github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanTransformBuffer struct {
	StagingBuffer sealBuffer.VulkanBuffer
	Device        *logical.VulkanLogicalDevice
	Options       VulkanTransformBufferOptions
}

func NewVulkanTransformBuffer(pipeline *sealPipeline.VulkanPipeline, options VulkanTransformBufferOptions) (VulkanTransformBuffer, error) {
	var err error
	transformBuffer := VulkanTransformBuffer{
		Device:  pipeline.Device,
		Options: options,
	}

	if transformBuffer.StagingBuffer, err = sealBuffer.CreateVulkanBufferWithContainer(&pipeline.BufferContainer, sealBuffer.VulkanBufferOptionsData{Size: GetVulkanTransformBufferOptionsSize(&options), Usage: vulkan.BufferUsageFlags(vulkan.BufferUsageTransferSrcBit | vulkan.BufferUsageUniformBufferBit), SharingMode: vulkan.SharingModeExclusive, Flags: vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit | vulkan.MemoryPropertyHostCoherentBit)}); err != nil {
		return transformBuffer, err
	}
	if err := CopyVulkanTransformBuffer(&transformBuffer); err != nil {
		return transformBuffer, err
	}

	return transformBuffer, nil
}

func CopyVulkanTransformBuffer(buffer *VulkanTransformBuffer) error {
	var vulkanUniformData unsafe.Pointer
	if res := vulkan.MapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory, GetVulkanTransformBufferOptionsUniformsOffset(&buffer.Options), GetVulkanTransformBufferOptionsUniformsSize(&buffer.Options), 0, &vulkanUniformData); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
		return vulkan.Error(res)
	}
	vulkanUniformBuffer := unsafe.Slice((*uniform.VulkanUniform)(vulkanUniformData), len(buffer.Options.Uniforms))
	copy(vulkanUniformBuffer, buffer.Options.Uniforms)
	vulkan.UnmapMemory(buffer.Device.Handle, buffer.StagingBuffer.Memory)

	return nil
}

func FreeVulkanTransformBuffer(buffer *VulkanTransformBuffer) error {
	return sealBuffer.FreeVulkanBuffer(&buffer.StagingBuffer)
}
