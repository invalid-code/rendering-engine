package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

var vertices = []float32{
	-0.5, 0.0, 0.0,
	-0.5, 0.5, 0.0,
	0.5, 0.0, 0.0,
}

const (
	vShaderSrc = `#version 330

layout (location = 0) in vec3 position;

void main() {
	gl_Position = vec4(position.x, position.y, position.z , 1.0);
}` + "\x00"
	fShaderSrc = `#version 330

layout (location = 0) out vec4 color;

void main() {
	color = vec4(0.0, 0.0, 1.0, 1.0);
}` + "\x00"
)

func compileShader(shaderSrc string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	cstring, free := gl.Strs(shaderSrc)
	gl.ShaderSource(shader, 1, cstring, nil)
	free()
	gl.CompileShader(shader)

	var compileStatus int32
	var shaderLog uint8 = 255
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &compileStatus)
	if compileStatus != gl.TRUE {
		gl.GetShaderInfoLog(shader, 255, nil, &shaderLog)
		log.Fatalf("\n%v\n", gl.GoStr(&shaderLog))
	}

	return shader
}

func createProgram() uint32 {
	vShader := compileShader(vShaderSrc, gl.VERTEX_SHADER)
	fShader := compileShader(fShaderSrc, gl.FRAGMENT_SHADER)
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	return program
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(600, 400, "heat rendering engine", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

func initVao() uint32 {
	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 4*3, uintptr(0))
	gl.EnableVertexAttribArray(0)
	return vao
}

func draw(vao uint32, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func redraw(window *glfw.Window) {
	window.SwapBuffers()
	glfw.PollEvents()
}

func handleInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func main() {
	window := initGlfw()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	program := createProgram()
	vao := initVao()

	for !window.ShouldClose() {
		handleInput(window)

		draw(vao, program)

		redraw(window)
	}
}
