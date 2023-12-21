package buffer

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	sealUniform "github.com/LamkasDev/seal/cmd/engine/vulkan/uniform"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/vulkan-go/vulkan"
)

type VulkanMeshBufferOptions struct {
	Vertices []vertex.VulkanVertex
	Indices  []uint16
	Uniforms []sealUniform.VulkanUniform
}

func NewVulkanMeshBufferOptions(vertices []vertex.VulkanVertex, indices []uint16, uniform sealUniform.VulkanUniform) VulkanMeshBufferOptions {
	data := VulkanMeshBufferOptions{
		Vertices: vertices,
		Indices:  indices,
		Uniforms: make([]sealUniform.VulkanUniform, commonPipeline.MaxFramesInFlight),
	}

	return data
}

func GetVulkanMeshBufferOptionsVerticesOffset(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return 0
}

func GetVulkanMeshBufferOptionsVerticesSize(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return vulkan.DeviceSize(int(vertex.VulkanVertexSize) * len(options.Vertices))
}

func GetVulkanMeshBufferOptionsIndicesOffset(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return GetVulkanMeshBufferOptionsVerticesOffset(options) + GetVulkanMeshBufferOptionsVerticesSize(options)
}

func GetVulkanMeshBufferOptionsIndicesSize(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return vulkan.DeviceSize(2 * len(options.Indices))
}

func GetVulkanMeshBufferOptionsUniformsOffset(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return GetVulkanMeshBufferOptionsIndicesOffset(options) + GetVulkanMeshBufferOptionsIndicesSize(options)
}

func GetVulkanMeshBufferOptionsUniformsSize(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return vulkan.DeviceSize(int(sealUniform.VulkanUniformSize) * len(options.Uniforms))
}

func GetVulkanMeshBufferOptionsSize(options *VulkanMeshBufferOptions) vulkan.DeviceSize {
	return GetVulkanMeshBufferOptionsVerticesSize(options) + GetVulkanMeshBufferOptionsIndicesSize(options) + GetVulkanMeshBufferOptionsUniformsSize(options)
}
