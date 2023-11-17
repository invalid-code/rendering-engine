package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	VERTICES = []float32{
		-0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0,  // bottom-left
		 0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  1.0, 0.0,  // bottom-right
		-0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0,  // top-left
		 0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  1.0, 1.0,  // top-right
	}
	INDICES = []uint32{
		0, 1, 2,  // first triangle
		1, 3, 2,  // second triangle
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
	var shaderLog uint8
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
	gl.UseProgram(program)
	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)
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

func initVao() (uint32, uint32, uint32) {
	var vbo, vao, ebo, texture uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.GenVertexArrays(1, &vao)
	gl.GenTextures(gl.TEXTURE_2D, &texture)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	f, err := os.Open("assets/texture/Brickwall2_Texture.png")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BufferData(gl.ARRAY_BUFFER, len(VERTICES)*4, gl.Ptr(VERTICES), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(INDICES)*4, gl.Ptr(INDICES), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(8*4), uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(8*4), uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, int32(8*4), uintptr(6*4))
	gl.EnableVertexAttribArray(2)
	return vao, ebo, texture
}

func drawBuffer(vao uint32, program uint32, ebo uint32, texture uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

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
	vao, ebo, texture := initVao()

	translate := mgl32.Rotate2D(45.0).Mat4()
	translateLocation := gl.GetUniformLocation(program, gl.Str("translate\x00"))
	gl.UniformMatrix4fv(translateLocation, 1, false, &translate[0])
	// mgl32.Rotate2D(45).Mat4()
	fmt.Println(translate)

	for !window.ShouldClose() {
		handleInput(window)

		drawBuffer(vao, program, ebo, texture)

		redraw(window)
	}
}
