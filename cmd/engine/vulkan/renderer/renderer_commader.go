package renderer

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
)

type VulkanRendererCommander struct {
	Device          *logical.VulkanLogicalDevice
	ShaderContainer *shader.VulkanShaderContainer
	Pool            command.VulkanCommandPool

	StagingBuffer          sealBuffer.VulkanCommandBuffer
	CommandBuffers         []sealBuffer.VulkanCommandBuffer
	CurrentCommandBuffer   *sealBuffer.VulkanCommandBuffer
	AbstractCommandBuffers map[string]sealBuffer.VulkanAbstractBuffer
}

func NewVulkanRendererCommander(device *logical.VulkanLogicalDevice, shaderContainer *shader.VulkanShaderContainer) (VulkanRendererCommander, error) {
	var err error
	commander := VulkanRendererCommander{
		Device:                 device,
		ShaderContainer:        shaderContainer,
		CommandBuffers:         make([]sealBuffer.VulkanCommandBuffer, commonPipeline.MaxFramesInFlight),
		AbstractCommandBuffers: map[string]sealBuffer.VulkanAbstractBuffer{},
	}

	if commander.Pool, err = command.NewVulkanCommandPool(device, uint32(device.Physical.Capabilities.Queue.GraphicsIndex)); err != nil {
		return commander, err
	}
	if commander.StagingBuffer, err = sealBuffer.NewVulkanCommandBuffer(device, &commander.Pool); err != nil {
		return commander, err
	}
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if commander.CommandBuffers[i], err = sealBuffer.NewVulkanCommandBuffer(device, &commander.Pool); err != nil {
			return commander, err
		}
	}
	for key := range shaderContainer.Shaders {
		commander.AbstractCommandBuffers[key] = sealBuffer.NewVulkanAbstractBuffer()
	}

	return commander, nil
}

func RecordVulkanRendererCommanderCommands(commander *VulkanRendererCommander, shader string, action sealBuffer.VulkanAbstractBufferAction) {
	buffer := commander.AbstractCommandBuffers[shader]
	buffer.Actions = append(buffer.Actions, action)
	commander.AbstractCommandBuffers[shader] = buffer
}

func RunVulkanRendererCommanderCommands(commander *VulkanRendererCommander, shader string) {
	for _, action := range commander.AbstractCommandBuffers[shader].Actions {
		action()
	}
}

func ResetVulkanRendererCommanderCommands(commander *VulkanRendererCommander, shader string) {
	buffer := commander.AbstractCommandBuffers[shader]
	buffer.Actions = []sealBuffer.VulkanAbstractBufferAction{}
	commander.AbstractCommandBuffers[shader] = buffer
}

func FreeVulkanRendererCommander(commander *VulkanRendererCommander) error {
	for i := 0; i < commonPipeline.MaxFramesInFlight; i++ {
		if err := sealBuffer.FreeVulkanCommandBuffer(&commander.CommandBuffers[i]); err != nil {
			return err
		}
	}
	if err := sealBuffer.FreeVulkanCommandBuffer(&commander.StagingBuffer); err != nil {
		return err
	}

	return command.FreeVulkanCommandPool(&commander.Pool)
}
