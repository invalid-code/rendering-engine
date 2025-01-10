package main

import (
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	vertices = []float32{
		// left top front(0)
		-0.5,  0.5,  0.5,
		// right top front(1)
		 0.5,  0.5,  0.5,
		// left bottom front(2)
		-0.5, -0.5,  0.5,
		// right bottom front(3)
		 0.5, -0.5,  0.5,
		// left top back(4)
		-0.5,  0.5, -0.5,
		// right top back(5)
		 0.5,  0.5, -0.5,
		// left bottom back(6)
		-0.5, -0.5, -0.5,
		// right bottom back(7)
		 0.5, -0.5, -0.5,
	}
	light_vertices = []float32{
		
	}
	indices = []uint32{
		// front
		2, 3, 1,
		1, 0, 2,
		// top
		0, 1, 5,
		5, 4, 0,
		// right
		3, 7, 5,
		5, 1, 3,
		// back
		7, 6, 4,
		4, 5, 7,
		// bottom
		7, 6, 2,
		2, 3, 7,
		// left
		6, 2, 0,
		0, 4, 6,
	}
	cubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}

	deltaTime float32 = 0.0
	lastFrame float32 = 0.0

	yaw   float32 = -90.0
	pitch float32 = 0.0

	lastX      float32 = WIDTH / 2
	lastY      float32 = HEIGHT / 2
	firstMouse         = true

	fov float32 = 45.0
)

const (
	WIDTH  = 600
	HEIGHT = 400
)

func init() {
	runtime.LockOSThread()
}

type Program struct {
	vShader uint32
	fShader uint32
	program uint32
}

func newShaderProgram(vShaderPath string, fShaderPath string) *Program {
	program := new(Program)
	program.vShader = program.compileShader(readFile(vShaderPath), gl.VERTEX_SHADER)
	program.fShader = program.compileShader(readFile(fShaderPath), gl.FRAGMENT_SHADER)
	program.program = gl.CreateProgram()
	gl.AttachShader(program.program, program.vShader)
	gl.AttachShader(program.program, program.fShader)
	gl.LinkProgram(program.program)
	return program
}

func (program *Program) deleteShader() {
	gl.DeleteShader(program.vShader)
	gl.DeleteShader(program.fShader)
}

func (program *Program) activate() {
	gl.UseProgram(program.program)
}

func (program *Program) setMat4(value mgl32.Mat4, name string) {
	location := gl.GetUniformLocation(program.program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func (program *Program) compileShader(shaderSrc string, shaderType uint32) uint32 {
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

type Window struct {
	window *glfw.Window
}

func newWindow() *Window {
	window := new(Window)
	glfwWindow, err := glfw.CreateWindow(WIDTH, HEIGHT, "heat rendering engine", nil, nil)
	if err != nil {
		panic(err)
	}
	glfwWindow.MakeContextCurrent()
	glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	glfwWindow.SetScrollCallback(func(window *glfw.Window, xOffset float64, yOffset float64) {
		fov -= float32(yOffset)
		if fov < 1.0 {
			fov = 1.0
		}
		if fov > 45.0 {
			fov = 45.0
		}
	})
	glfwWindow.SetCursorPosCallback(func(window *glfw.Window, xPos float64, yPos float64) {
		if firstMouse {
			lastX = float32(xPos)
			lastY = float32(yPos)
			firstMouse = false
		}
		xOffset := float32(xPos) - lastX
		yOffset := lastY - float32(yPos)
		lastX = float32(xPos)
		lastY = float32(yPos)

		var sensitivity float32 = 0.1
		xOffset *= sensitivity
		yOffset *= sensitivity

		yaw += float32(xOffset)
		pitch += float32(yOffset)

		if pitch > 89.0 {
			pitch = 89.0
		}
		if pitch < -89.0 {
			pitch = -89.0
		}

		direction := mgl32.Vec3{}
		direction[0] = float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch))))
		direction[1] = float32(math.Sin(float64(mgl32.DegToRad(pitch))))
		direction[2] = float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch))))
		cameraFront = cameraFront.Add(direction).Normalize()
	})
	window.window = glfwWindow
	return window
}

func (window *Window) mainLoop(f func()) {
	glfwWindow := window.getGlfwWindow()
	for !glfwWindow.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		cameraSpeed := 2.5 * deltaTime

		if glfwWindow.GetKey(glfw.KeyEscape) == glfw.Press {
			glfwWindow.SetShouldClose(true)
		}

		if glfwWindow.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}

		if glfwWindow.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}

		if glfwWindow.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		if glfwWindow.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		f()

		glfwWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func (window *Window) getGlfwWindow() *glfw.Window {
	return window.window
}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(file) + "\x00"
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)

	window := newWindow()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	program := newShaderProgram("assets/shader/basic_trianglev.glsl", "assets/shader/basic_trianglef.glsl")
	defer program.deleteShader()

	var verticesBufferObj, colorBufferObj, texCoordsBufferObj, vao, ebo, texture uint32
	gl.GenBuffers(1, &colorBufferObj)
	gl.GenBuffers(1, &texCoordsBufferObj)
	gl.GenBuffers(1, &ebo)
	
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	
	gl.GenTextures(gl.TEXTURE_2D, &texture)

	gl.BindBuffer(gl.ARRAY_BUFFER, verticesBufferObj)
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

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(8*4), uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(8*4), uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, int32(8*4), uintptr(6*4))
	gl.EnableVertexAttribArray(2)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	window.mainLoop(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.activate()
		gl.BindVertexArray(vao)

		projection := mgl32.Perspective(mgl32.DegToRad(fov), WIDTH/HEIGHT, 0.1, 100.0)
		program.setMat4(projection, "projection")

		view := mgl32.LookAt(cameraPos[0], cameraPos[1], cameraPos[2], cameraPos[0]+cameraFront[0], cameraPos[1]+cameraFront[1], cameraPos[2]+cameraFront[2], cameraUp[0], cameraUp[1], cameraUp[2])
		program.setMat4(view, "view")

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		
		for i := range cubePositions {
			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(cubePositions[i][0], cubePositions[i][1], cubePositions[i][2]))
			angle := float32(i * 20.0)
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{0.5, 1.0, 0.0}))
			program.setMat4(model, "model")

			gl.DrawElements(gl.TRIANGLES, 6*6, gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
		}
	})
}
