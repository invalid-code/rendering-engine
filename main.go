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
		// first square
		-0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		-0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
		// second square
		-0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		-0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
		// third square
		0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
		// fourth square
		-0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		-0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		-0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		-0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
		// fifth square
		-0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		-0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
		// sixth square
		-0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
		-0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
		0.5, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
		0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
	}
	indices = []uint32{
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
	cubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -15.0},
		mgl32.Vec3{-1.5, -2.2, -2.5},
		mgl32.Vec3{-3.8, -2.0, -12.3},
		mgl32.Vec3{2.4, -0.4, -3.5},
		mgl32.Vec3{-1.7, 3.0, -7.5},
		mgl32.Vec3{1.3, -2.0, -2.5},
		mgl32.Vec3{1.5, 2.0, -2.5},
		mgl32.Vec3{1.5, 0.2, -1.5},
		mgl32.Vec3{-1.3, 1.0, -1.5},
	}

	cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}

	deltaTime float32 = 0.0
	lastFrame float32 = 0.0

	yaw float32 = -90.0
	pitch float32 = 0.0

	lastX float32 = WIDTH / 2
	lastY float32 = HEIGHT / 2
	firstMouse = true

	fov float32 = 45.0
)

const (
	WIDTH  = 600
	HEIGHT = 400
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

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "heat rendering engine", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetScrollCallback(func(window *glfw.Window, xOffset float64, yOffset float64) {
		fov -= float32(yOffset)
		if fov < 1.0 {
			fov = 1.0
		}
		if fov > 45.0 {
			fov = 45.0
		}
	})
	window.SetCursorPosCallback(func(window *glfw.Window, xPos float64, yPos float64) {
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
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vShader := compileShader(readFile("assets/shader/basic_trianglev.glsl"), gl.VERTEX_SHADER)
	fShader := compileShader(readFile("assets/shader/basic_trianglef.glsl"), gl.FRAGMENT_SHADER)
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	defer gl.DeleteShader(vShader)
	defer gl.DeleteShader(fShader)

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

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(8*4), uintptr(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(8*4), uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, int32(8*4), uintptr(6*4))
	gl.EnableVertexAttribArray(2)

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		cameraSpeed := 2.5 * deltaTime

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		if window.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}

		if window.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}

		if window.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		if window.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		projection := mgl32.Perspective(mgl32.DegToRad(fov), WIDTH/HEIGHT, 0.1, 100.0)
		projectionLocation := gl.GetUniformLocation(program, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

		view := mgl32.LookAt(cameraPos[0], cameraPos[1], cameraPos[2], cameraPos[0]+cameraFront[0], cameraPos[1]+cameraFront[1], cameraPos[2]+cameraFront[2], cameraUp[0], cameraUp[1], cameraUp[2])
		viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		for i := range cubePositions {
			model := mgl32.Ident4()
			modelLocation := gl.GetUniformLocation(program, gl.Str("model\x00"))
			model = model.Mul4(mgl32.Translate3D(cubePositions[i][0], cubePositions[i][1], cubePositions[i][2]))
			angle := float32(i * 20.0)
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{0.5, 1.0, 0.0}))
			gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])

			gl.DrawElements(gl.TRIANGLES, 6*6, gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
