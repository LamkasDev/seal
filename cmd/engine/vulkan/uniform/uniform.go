package uniform

import (
	"unsafe"

	"github.com/EngoEngine/glm"
	"github.com/vulkan-go/vulkan"
)

type VulkanUniform struct {
	Model      glm.Mat4
	View       glm.Mat4
	Projection glm.Mat4
}

const VulkanUniformSize = unsafe.Sizeof(VulkanUniform{})
const VulkanUniformModelOffset = unsafe.Offsetof(VulkanUniform{}.Model)
const VulkanUniformViewOffset = unsafe.Offsetof(VulkanUniform{}.View)
const VulkanUniformProjectionOffset = unsafe.Offsetof(VulkanUniform{}.Projection)

func NewVulkanUniform(extent vulkan.Extent2D, position glm.Vec3, rotation float32) VulkanUniform {
	uniform := VulkanUniform{
		Model:      glm.HomogRotate3D(glm.DegToRad(rotation), &glm.Vec3{0, 0, 1}),
		View:       glm.LookAt(2, 2, 2, position.X(), position.Y(), position.Z(), 0, 0, 1),
		Projection: glm.Perspective(glm.DegToRad(45), float32(extent.Width)/float32(extent.Height), 0.1, 10),
	}
	uniform.Projection[5] *= -1

	return uniform
}
