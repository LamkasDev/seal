package renderer

import (
	commonPipeline "github.com/LamkasDev/seal/cmd/common/pipeline"
	sealBuffer "github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/command"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pass"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/vulkan-go/vulkan"
)

type VulkanRendererCommander struct {
	Device              *logical.VulkanLogicalDevice
	ShaderContainer     *shader.VulkanShaderContainer
	RenderPassContainer pass.VulkanRenderPassContainer
	Pool                command.VulkanCommandPool

	StagingBuffer        sealBuffer.VulkanCommandBuffer
	CommandBuffers       []sealBuffer.VulkanCommandBuffer
	CurrentCommandBuffer *sealBuffer.VulkanCommandBuffer
}

func NewVulkanRendererCommander(device *logical.VulkanLogicalDevice, format vulkan.Format, shaders []string) (VulkanRendererCommander, error) {
	var err error
	commander := VulkanRendererCommander{
		Device:         device,
		CommandBuffers: make([]sealBuffer.VulkanCommandBuffer, commonPipeline.MaxFramesInFlight),
	}

	if commander.RenderPassContainer, err = pass.NewVulkanRenderPassContainer(device, format, shaders); err != nil {
		return commander, err
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

	return commander, nil
}

func RecordVulkanRendererCommanderCommands(commander *VulkanRendererCommander, layer uint8, shader string, action sealBuffer.VulkanAbstractBufferAction) {
	buffer := commander.RenderPassContainer.Passes[layer].AbstractCommandBuffers[shader]
	buffer.Actions = append(buffer.Actions, action)
	commander.RenderPassContainer.Passes[layer].AbstractCommandBuffers[shader] = buffer
}

func RunVulkanRendererCommanderCommands(commander *VulkanRendererCommander, layer uint8, shader string) {
	for _, action := range commander.RenderPassContainer.Passes[layer].AbstractCommandBuffers[shader].Actions {
		action()
	}
}

func ResetVulkanRendererCommanderCommands(commander *VulkanRendererCommander, layer uint8, shader string) {
	buffer := commander.RenderPassContainer.Passes[layer].AbstractCommandBuffers[shader]
	buffer.Actions = []sealBuffer.VulkanAbstractBufferAction{}
	commander.RenderPassContainer.Passes[layer].AbstractCommandBuffers[shader] = buffer
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
