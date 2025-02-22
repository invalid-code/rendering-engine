package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WIDTH  = 685
	HEIGHT = 500
)

var (
	lightVertices = []float32{
		// left top front 0
		-0.5, 0.5, 0.5,
		// right top front 1
		0.5, 0.5, 0.5,
		// left bottom front 2
		-0.5, -0.5, 0.5,
		// right bottom front 3
		0.5, -0.5, 0.5,
		// left bottom back 4
		-0.5, 0.5, -0.5,
		// right bottom back 5
		0.5, 0.5, -0.5,
		// right top back 6
		-0.5, -0.5, -0.5,
		// left top back 7
		0.5, -0.5, -0.5,
	}
	indices = []int32{
		// front
		0, 2, 3,
		3, 1, 0,
		// back
		5, 7, 6,
		6, 4, 5,
		// right
		1, 3, 7,
		7, 5, 1,
		// left
		4, 6, 2,
		2, 0, 4,
		// top
		4, 0, 1,
		1, 5, 4,
		// bottom
		7, 3, 2,
		2, 6, 7,
	}
	vertices = []float32{
		// front
		lightVertices[indices[0]*3], lightVertices[(indices[0]*3)+1], lightVertices[(indices[0]*3)+2], 0.0, 0.0, 1.0, 0.0, 1.0,
		lightVertices[indices[1]*3], lightVertices[(indices[1]*3)+1], lightVertices[(indices[1]*3)+2], 0.0, 0.0, 1.0, 0.0, 0.0,
		lightVertices[indices[2]*3], lightVertices[(indices[2]*3)+1], lightVertices[(indices[2]*3)+2], 0.0, 0.0, 1.0, 1.0, 0.0,
		lightVertices[indices[3]*3], lightVertices[(indices[3]*3)+1], lightVertices[(indices[3]*3)+2], 0.0, 0.0, 1.0, 1.0, 0.0,
		lightVertices[indices[4]*3], lightVertices[(indices[4]*3)+1], lightVertices[(indices[4]*3)+2], 0.0, 0.0, 1.0, 1.0, 1.0,
		lightVertices[indices[5]*3], lightVertices[(indices[5]*3)+1], lightVertices[(indices[5]*3)+2], 0.0, 0.0, 1.0, 0.0, 1.0,
		// back
		lightVertices[indices[6]*3], lightVertices[(indices[6]*3)+1], lightVertices[(indices[6]*3)+2], 0.0, 0.0, -1.0, 0.0, 1.0,
		lightVertices[indices[7]*3], lightVertices[(indices[7]*3)+1], lightVertices[(indices[7]*3)+2], 0.0, 0.0, -1.0, 0.0, 0.0,
		lightVertices[indices[8]*3], lightVertices[(indices[8]*3)+1], lightVertices[(indices[8]*3)+2], 0.0, 0.0, -1.0, 1.0, 0.0,
		lightVertices[indices[9]*3], lightVertices[(indices[9]*3)+1], lightVertices[(indices[9]*3)+2], 0.0, 0.0, -1.0, 1.0, 0.0,
		lightVertices[indices[10]*3], lightVertices[(indices[10]*3)+1], lightVertices[(indices[10]*3)+2], 0.0, 0.0, -1.0, 1.0, 1.0,
		lightVertices[indices[11]*3], lightVertices[(indices[11]*3)+1], lightVertices[(indices[11]*3)+2], 0.0, 0.0, -1.0, 0.0, 1.0,
		// right
		lightVertices[indices[12]*3], lightVertices[(indices[12]*3)+1], lightVertices[(indices[12]*3)+2], 1.0, 0.0, 0.0, 0.0, 1.0,
		lightVertices[indices[13]*3], lightVertices[(indices[13]*3)+1], lightVertices[(indices[13]*3)+2], 1.0, 0.0, 0.0, 0.0, 0.0,
		lightVertices[indices[14]*3], lightVertices[(indices[14]*3)+1], lightVertices[(indices[14]*3)+2], 1.0, 0.0, 0.0, 1.0, 0.0,
		lightVertices[indices[15]*3], lightVertices[(indices[15]*3)+1], lightVertices[(indices[15]*3)+2], 1.0, 0.0, 0.0, 1.0, 0.0,
		lightVertices[indices[16]*3], lightVertices[(indices[16]*3)+1], lightVertices[(indices[16]*3)+2], 1.0, 0.0, 0.0, 1.0, 1.0,
		lightVertices[indices[17]*3], lightVertices[(indices[17]*3)+1], lightVertices[(indices[17]*3)+2], 1.0, 0.0, 0.0, 0.0, 1.0,
		// left
		lightVertices[indices[18]*3], lightVertices[(indices[18]*3)+1], lightVertices[(indices[18]*3)+2], -1.0, 0.0, 0.0, 0.0, 1.0,
		lightVertices[indices[19]*3], lightVertices[(indices[19]*3)+1], lightVertices[(indices[19]*3)+2], -1.0, 0.0, 0.0, 0.0, 0.0,
		lightVertices[indices[20]*3], lightVertices[(indices[20]*3)+1], lightVertices[(indices[20]*3)+2], -1.0, 0.0, 0.0, 1.0, 0.0,
		lightVertices[indices[21]*3], lightVertices[(indices[21]*3)+1], lightVertices[(indices[21]*3)+2], -1.0, 0.0, 0.0, 1.0, 0.0,
		lightVertices[indices[22]*3], lightVertices[(indices[22]*3)+1], lightVertices[(indices[22]*3)+2], -1.0, 0.0, 0.0, 1.0, 1.0,
		lightVertices[indices[23]*3], lightVertices[(indices[23]*3)+1], lightVertices[(indices[23]*3)+2], -1.0, 0.0, 0.0, 0.0, 1.0,
		// top
		lightVertices[indices[24]*3], lightVertices[(indices[24]*3)+1], lightVertices[(indices[24]*3)+2], 0.0, 1.0, 0.0, 0.0, 1.0,
		lightVertices[indices[25]*3], lightVertices[(indices[25]*3)+1], lightVertices[(indices[25]*3)+2], 0.0, 1.0, 0.0, 0.0, 0.0,
		lightVertices[indices[26]*3], lightVertices[(indices[26]*3)+1], lightVertices[(indices[26]*3)+2], 0.0, 1.0, 0.0, 1.0, 0.0,
		lightVertices[indices[27]*3], lightVertices[(indices[27]*3)+1], lightVertices[(indices[27]*3)+2], 0.0, 1.0, 0.0, 1.0, 0.0,
		lightVertices[indices[28]*3], lightVertices[(indices[28]*3)+1], lightVertices[(indices[28]*3)+2], 0.0, 1.0, 0.0, 1.0, 1.0,
		lightVertices[indices[29]*3], lightVertices[(indices[29]*3)+1], lightVertices[(indices[29]*3)+2], 0.0, 1.0, 0.0, 0.0, 1.0,
		// bottom
		lightVertices[indices[30]*3], lightVertices[(indices[30]*3)+1], lightVertices[(indices[30]*3)+2], 0.0, -1.0, 0.0, 0.0, 1.0,
		lightVertices[indices[31]*3], lightVertices[(indices[31]*3)+1], lightVertices[(indices[31]*3)+2], 0.0, -1.0, 0.0, 0.0, 0.0,
		lightVertices[indices[32]*3], lightVertices[(indices[32]*3)+1], lightVertices[(indices[32]*3)+2], 0.0, -1.0, 0.0, 1.0, 0.0,
		lightVertices[indices[33]*3], lightVertices[(indices[33]*3)+1], lightVertices[(indices[33]*3)+2], 0.0, -1.0, 0.0, 1.0, 0.0,
		lightVertices[indices[34]*3], lightVertices[(indices[34]*3)+1], lightVertices[(indices[34]*3)+2], 0.0, -1.0, 0.0, 1.0, 1.0,
		lightVertices[indices[35]*3], lightVertices[(indices[35]*3)+1], lightVertices[(indices[35]*3)+2], 0.0, -1.0, 0.0, 0.0, 1.0,
	}
	pointLightsPositions         = []mgl32.Vec3{mgl32.Vec3(randFloats(-5.0, 5.0, 3)), mgl32.Vec3(randFloats(-5.0, 5.0, 3)), mgl32.Vec3(randFloats(-5.0, 5.0, 3)), mgl32.Vec3(randFloats(-5.0, 5.0, 3))}
	deltaTime            float32 = 0.0
	lastFrame            float32 = 0.0
	lastX                float32 = 400
	lastY                float32 = 300
	lightColor                   = mgl32.Vec3{1.0, 1.0, 1.0}
	ambient                      = mgl32.Vec3{0.2, 0.2, 0.2}
	diffuse                      = mgl32.Vec3{0.7, 0.7, 0.7}
	specular                     = mgl32.Vec3{1.0, 1.0, 1.0}
)

type Camera struct {
	pos        mgl32.Vec3
	direction  mgl32.Vec3
	up         mgl32.Vec3
	fov        float32
	yaw        float32
	pitch      float32
	speed      float32
	firstMouse bool
}

func newCamera() Camera {
	camera := Camera{
		pos:        mgl32.Vec3{0, 0, 3.0},
		direction:  mgl32.Vec3{0, 0, -1.0},
		up:         mgl32.Vec3{0, 1.0, 0},
		fov:        45.0,
		yaw:        -90.0,
		pitch:      0,
		speed:      0,
		firstMouse: true,
	}
	return camera
}

func (camera *Camera) zoom(zoomVal float32) {
	camera.fov -= zoomVal
	if camera.fov < 1.0 {
		camera.fov = 1.0
	}
	if camera.fov > 45.0 {
		camera.fov = 45.0
	}
}

func (camera *Camera) cameraTarget() mgl32.Vec3 {
	return camera.pos.Add(camera.direction)
}

func (camera *Camera) calculateDirection(xPos float32, yPos float32) {
	if camera.firstMouse {
		lastX = xPos
		lastY = yPos
		camera.firstMouse = false
	}
	xOffset, yOffset := xPos-lastX, lastY-yPos
	lastX = xPos
	lastY = yPos

	const sensitivity float32 = 0.1
	xOffset *= sensitivity
	yOffset *= sensitivity

	camera.yaw += xOffset
	camera.pitch += yOffset

	if camera.pitch > 89.0 {
		camera.pitch = 89.0
	}
	if camera.pitch < -89.0 {
		camera.pitch = -89.0
	}
	var direction mgl32.Vec3
	direction[0] = float32(math.Cos(float64(mgl32.DegToRad(camera.yaw))) * math.Cos(float64(mgl32.DegToRad(camera.pitch))))
	direction[1] = float32(math.Sin(float64(mgl32.DegToRad(camera.pitch))))
	direction[2] = float32(math.Sin(float64(mgl32.DegToRad(camera.yaw))) * math.Cos(float64(mgl32.DegToRad(camera.pitch))))
	camera.direction = direction.Normalize()
}

func (camera *Camera) calculateSpeed() {
	camera.speed = 2.5 * deltaTime
}

func (camera *Camera) moveForward() {
	camera.pos = camera.pos.Add(camera.direction.Mul(camera.speed))
}

func (camera *Camera) moveBackward() {
	camera.pos = camera.pos.Sub(camera.direction.Mul(camera.speed))
}

func (camera *Camera) moveRight() {
	camera.pos = camera.pos.Add(camera.direction.Cross(camera.up).Normalize().Mul(camera.speed))
}

func (camera *Camera) moveLeft() {
	camera.pos = camera.pos.Sub(camera.direction.Cross(camera.up).Normalize().Mul(camera.speed))
}

func (camera *Camera) moveUp() {
	camera.pos = camera.pos.Add(camera.up.Mul(camera.speed))
}

func (camera *Camera) moveDown() {
	camera.pos = camera.pos.Sub(camera.up.Mul(camera.speed))
}

type ShaderProgram struct {
	program uint32
}

func newShaderProgram(vertShaderPath string, fragShaderPath string) ShaderProgram {
	shaderPaths := []string{vertShaderPath, fragShaderPath}
	shaderType := []uint32{gl.VERTEX_SHADER, gl.FRAGMENT_SHADER}
	compiledShader := [2]uint32{}
	for i := 0; i < len(shaderPaths); i++ {
		compiledShader[i] = compileShader(shaderPaths[i], shaderType[i])
	}
	shaderProgram := ShaderProgram{
		program: gl.CreateProgram(),
	}
	for i := 0; i < len(compiledShader); i++ {
		gl.AttachShader(shaderProgram.program, compiledShader[i])
	}
	gl.LinkProgram(shaderProgram.program)
	var success int32
	infoLog := gl.Str(string(make([]byte, 512)) + "\x00")
	gl.GetProgramiv(shaderProgram.program, gl.LINK_STATUS, &success)
	if success != 1 {
		gl.GetProgramInfoLog(shaderProgram.program, 512, nil, infoLog)
		log.Fatalln(gl.GoStr(infoLog))
	}
	for i := 0; i < len(compiledShader); i++ {
		gl.DeleteShader(compiledShader[i])
	}
	return shaderProgram
}

func compileShader(path string, shaderType uint32) uint32 {
	var shader uint32
	shaderSrc := readFile(path)
	cStr, freeFn := gl.Strs(shaderSrc + "\x00")
	defer freeFn()
	shader = gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, cStr, nil)
	gl.CompileShader(shader)
	var success int32
	infoLog := gl.Str(string(make([]byte, 512)) + "\x00")
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success != 1 {
		gl.GetShaderInfoLog(shader, 512, nil, infoLog)
		log.Fatalln(path, gl.GoStr(infoLog))
	}
	return shader
}

func (shaderProgram *ShaderProgram) setMat4(name string, uniformData mgl32.Mat4) {
	gl.UseProgram(shaderProgram.program)
	uniformName := gl.Str(name + "\x00")
	uniformLoc := gl.GetUniformLocation(shaderProgram.program, uniformName)
	gl.UniformMatrix4fv(uniformLoc, 1, false, &uniformData[0])
	gl.UseProgram(0)
}

func (shaderProgram *ShaderProgram) setVec3(name string, uniformData mgl32.Vec3) {
	gl.UseProgram(shaderProgram.program)
	uniformName := gl.Str(name + "\x00")
	uniformLoc := gl.GetUniformLocation(shaderProgram.program, uniformName)
	gl.Uniform3fv(uniformLoc, 1, &uniformData[0])
	gl.UseProgram(0)
}

func (shaderProgram *ShaderProgram) setFloat(name string, uniformData float32) {
	gl.UseProgram(shaderProgram.program)
	uniformName := gl.Str(name + "\x00")
	uniformLoc := gl.GetUniformLocation(shaderProgram.program, uniformName)
	gl.Uniform1f(uniformLoc, uniformData)
	gl.UseProgram(0)
}

func (shaderProgram *ShaderProgram) setInt(name string, uniformData int32) {
	gl.UseProgram(shaderProgram.program)
	uniformName := gl.Str(name + "\x00")
	uniformLoc := gl.GetUniformLocation(shaderProgram.program, uniformName)
	gl.Uniform1i(uniformLoc, uniformData)
	gl.UseProgram(0)
}

func (shaderProgram *ShaderProgram) activate() {
	gl.UseProgram(shaderProgram.program)
}

func init() {
	runtime.LockOSThread()
}

func loadTexture(imgPath string) uint32 {
	imgBytes := []byte(readFile(imgPath))
	imgReader := bytes.NewReader(imgBytes)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		log.Fatalln(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	return texture
}

func createVAOwithEBO() uint32 {
	var vbo, ebo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(lightVertices)*int(unsafe.Sizeof(lightVertices[0])), gl.Ptr(lightVertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(indices[0])), gl.Ptr(indices), gl.STATIC_DRAW)

	floatSize := unsafe.Sizeof(float32(1.0))
	stride := int32(3 * floatSize)
	offset := uintptr(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, offset)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	return vao
}

func createVAO() uint32 {
	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(vertices[0])), gl.Ptr(vertices), gl.STATIC_DRAW)

	floatSize := unsafe.Sizeof(float32(1.0))
	stride := int32(8 * floatSize)
	offset := uintptr(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, offset)
	gl.EnableVertexAttribArray(0)

	offset += 3 * floatSize
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, stride, offset)
	gl.EnableVertexAttribArray(1)

	offset += 3 * floatSize
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, stride, offset)
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)
	return vao
}

func processInput(window *glfw.Window, camera *Camera) {
	camera.calculateSpeed()
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.moveForward()
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.moveBackward()
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.moveLeft()
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.moveRight()
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		camera.moveUp()
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		camera.moveDown()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	err := glfw.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "heat rendering engine", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		log.Fatalln(err)
	}
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	camera := newCamera()
	window.SetCursorPosCallback(func(w *glfw.Window, xPos float64, yPos float64) {
		camera.calculateDirection(float32(xPos), float32(yPos))
	})
	window.SetScrollCallback(func(w *glfw.Window, xOff float64, yOff float64) {
		camera.zoom(float32(yOff))
	})
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.CULL_FACE)
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	lightVao := createVAOwithEBO()
	vao := createVAO()
	shaderProgram := newShaderProgram("assets/shader/mainCubev.glsl", "assets/shader/mainCubef.glsl")
	lightShaderProgram := newShaderProgram("assets/shader/lightCubev.glsl", "assets/shader/lightCubef.glsl")

	shaderProgram.setInt("material.diffuse", 0)
	shaderProgram.setInt("material.specular", 1)
	shaderProgram.setVec3("material.specular", mgl32.Vec3{0.5, 0.5, 0.5})
	shaderProgram.setFloat("material.shinniness", 32.0)
	shaderProgram.setVec3("dirLight.direction", mgl32.Vec3{0.0, -1.0, 0.0})
	shaderProgram.setVec3("dirLight.ambient", ambient)
	shaderProgram.setVec3("dirLight.diffuse", diffuse)
	shaderProgram.setVec3("dirLight.specular", specular)
	for i, pointLightPos := range pointLightsPositions {
		shaderProgram.setVec3(fmt.Sprintf("pointLights[%d].pos", i), pointLightPos)
		shaderProgram.setVec3(fmt.Sprintf("pointLights[%d].ambient", i), ambient)
		shaderProgram.setVec3(fmt.Sprintf("pointLights[%d].diffuse", i), diffuse)
		shaderProgram.setVec3(fmt.Sprintf("pointLights[%d].specular", i), specular)
		shaderProgram.setFloat(fmt.Sprintf("pointLights[%d].constant", i), 1.0)
		shaderProgram.setFloat(fmt.Sprintf("pointLights[%d].linear", i), 0.045)
		shaderProgram.setFloat(fmt.Sprintf("pointLights[%d].quadratic", i), 0.0075)
	}
	shaderProgram.setVec3("spotLight.ambient", ambient)
	shaderProgram.setVec3("spotLight.diffuse", diffuse)
	shaderProgram.setVec3("spotLight.specular", specular)
	shaderProgram.setFloat("spotLight.cutOff", float32(math.Cos(float64(mgl32.DegToRad(9.5)))))
	shaderProgram.setFloat("spotLight.outerCutOff", float32(math.Cos(float64(mgl32.DegToRad(12.5)))))
	shaderProgram.setFloat("spotLight.constant", 1.0)
	shaderProgram.setFloat("spotLight.linear", 0.045)
	shaderProgram.setFloat("spotLight.quadratic", 0.0075)

	texture := loadTexture("assets/texture/Brick_01-256x256.png")
	specTexture := loadTexture("assets/texture/spec_img.png")

	pos := make([][]float32, 10)
	rotAngle := make([][]float32, 10)
	for i := range pos {
		pos[i] = randFloats(-3.0, 3.0, 3)
		rotAngle[i] = randFloats(1.0, 45.0, 3)
	}
	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame
		processInput(window, &camera)

		gl.ClearColor(0.04, 0.12, 0.19, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		view := mgl32.LookAtV(camera.pos, camera.cameraTarget(), camera.up)
		projection := mgl32.Perspective(mgl32.DegToRad(camera.fov), float32(WIDTH)/float32(HEIGHT), 0.1, 100.0)

		lightShaderProgram.setMat4("view", view)
		lightShaderProgram.setMat4("projection", projection)
		lightShaderProgram.setVec3("lightColor", lightColor)
		gl.BindVertexArray(lightVao)
		for _, pointLightPosition := range pointLightsPositions {
			lightModel := mgl32.Ident4()
			lightModel = lightModel.Mul4(mgl32.Translate3D(pointLightPosition[0], pointLightPosition[1], pointLightPosition[2]))
			lightShaderProgram.setMat4("model", lightModel)
			lightShaderProgram.activate()
			gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
		}

		shaderProgram.setMat4("view", view)
		shaderProgram.setMat4("projection", projection)
		shaderProgram.setVec3("viewPos", camera.pos)
		shaderProgram.setVec3("spotLight.pos", camera.pos)
		shaderProgram.setVec3("spotLight.direction", camera.direction)
		gl.BindVertexArray(vao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, specTexture)
		for i := 0; i < 10; i++ {
			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Rotate3DX(rotAngle[i][0]).Mat4())
			model = model.Mul4(mgl32.Rotate3DY(rotAngle[i][1]).Mat4())
			model = model.Mul4(mgl32.Rotate3DZ(rotAngle[i][2]).Mat4())
			model = model.Mul4(mgl32.Translate3D(pos[0][0], pos[0][1], pos[0][2]))
			shaderProgram.setMat4("model", model)
			shaderProgram.activate()
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
