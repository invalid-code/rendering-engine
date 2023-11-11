package main

import (
	"log"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	floatSize = 4
	attribs   = 3
)

var (
	vertices = []float32{
		-0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, 0.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.0, 0.0, 0.0, 1.0,
		// 0.5, 0.5, 0.0, 0.0, 0.0, 1.0,
	}
	indices = []uint32{
		0, 1, 2,
		// 1, 3, 2,
	}
)

func init() {
	runtime.LockOSThread()
}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(file) + "\x00"
}

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
	vShader := compileShader(readFile("assets/shader/basic_trianglev.glsl"), gl.VERTEX_SHADER)
	fShader := compileShader(readFile("assets/shader/basic_trianglef.glsl"), gl.FRAGMENT_SHADER)
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

func initVao() (uint32, uint32) {
	var vbo, vao, ebo uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*floatSize, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, attribs, gl.FLOAT, false, int32(floatSize*len(vertices)/attribs), uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, attribs, gl.FLOAT, false, int32(floatSize*len(vertices)/attribs), uintptr(attribs*floatSize))
	gl.EnableVertexAttribArray(1)
	return vao, ebo
}

func draw(vao uint32, program uint32, ebo uint32) {
	time := glfw.GetTime()
	transperancyUniformLocation := gl.GetUniformLocation(program, gl.Str("transperancy\x00"))
	transperancyUniform := math.Sin(time) / 0.5

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)

	gl.Uniform1f(transperancyUniformLocation, float32(transperancyUniform))

	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
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
	defer glfw.Terminate()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	program := createProgram()
	vao, ebo := initVao()

	for !window.ShouldClose() {
		handleInput(window)

		draw(vao, program, ebo)

		redraw(window)
	}
}
