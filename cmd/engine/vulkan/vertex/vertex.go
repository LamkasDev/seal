package vertex

import (
	"unsafe"

	"github.com/EngoEngine/glm"
)

type VulkanVertex struct {
	Position glm.Vec2
	Color    glm.Vec3
}

const VulkanVertexSize = unsafe.Sizeof(VulkanVertex{})
const VulkanVertexPositionOffset = unsafe.Offsetof(VulkanVertex{}.Position)
const VulkanVertexColorOffset = unsafe.Offsetof(VulkanVertex{}.Color)

var DefaultVertices = []VulkanVertex{
	{Position: glm.Vec2{-0.5, -0.5}, Color: glm.Vec3{1, 0, 0}},
	{Position: glm.Vec2{0.5, -0.5}, Color: glm.Vec3{0, 1, 0}},
	{Position: glm.Vec2{0.5, 0.5}, Color: glm.Vec3{0, 0, 1}},
	{Position: glm.Vec2{-0.5, 0.5}, Color: glm.Vec3{1, 1, 1}},
}
var DefaultVerticesIndex = []uint16{
	0, 1, 2, 2, 3, 0,
}
