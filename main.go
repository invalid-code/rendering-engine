package main

import (
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
		// first square
		-0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		 0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		-0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		 0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
		// second square
		-0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		 0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		-0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		 0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
		// third square
		 0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		 0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		 0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		 0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
		// fourth square
		-0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		-0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		-0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		-0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
		// fifth square
		-0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		-0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		 0.5, 0.5, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		 0.5, 0.5, 1.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
		// sixth square
		-0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 0.0, // bottom-left
		-0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  1.0, 0.0, // bottom-right
		 0.5, 0.0, 0.0,  0.0, 0.0, 1.0,  0.0, 1.0, // top-left
		 0.5, 0.0, 1.0,  0.0, 0.0, 1.0,  1.0, 1.0, // top-right
	}
	INDICES = []uint32{
		// first square
		0, 1, 2, // first triangle
		1, 2, 3, // second triangle
		// second square
		4, 5, 6,
		5, 6, 7,
		// third square
		8, 9, 10,
		9, 10, 11,
		// fourth square
		12, 13, 14,
		13, 14, 15,
		// fifth square
		16, 17, 18,
		17, 18, 19,
		// sixth square
		20, 21, 22,
		21, 22, 23,
	}
	CUBE_POSITIONS = []mgl32.Vec3{
		mgl32.Vec3{0.0 ,0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
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
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawElements(gl.TRIANGLES, 6*6, gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
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

	projection := mgl32.Perspective(mgl32.DegToRad(90.0), 600/400, 0.1, 100.0)
	projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])
	
	view := mgl32.LookAt(0.0, 0.0, 2.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		handleInput(window)

		model := mgl32.Ident4()
		modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
		model = model.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()) * mgl32.DegToRad(50.0), mgl32.Vec3{0.5, 1.0, 0.0}))
		gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])
		
		drawBuffer(vao, program, ebo, texture)

		redraw(window)
	}
}
