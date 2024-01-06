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

func NewVulkanUniform(extent vulkan.Extent2D, camera glm.Vec3, position glm.Vec3, rotation float32) VulkanUniform {
	uniform := VulkanUniform{
		Model:      glm.Translate3D(position.X(), position.Y(), position.Z()),
		View:       glm.Translate3D(-camera.X(), -camera.Y(), -camera.Z()),
		Projection: glm.Perspective(glm.DegToRad(60), float32(extent.Width)/float32(extent.Height), 0.1, 10),
	}
	rotationMat := glm.HomogRotate3DY(glm.DegToRad(rotation))
	uniform.Model = rotationMat.Mul4(&uniform.Model)
	uniform.Projection[5] *= -1

	return uniform
}
