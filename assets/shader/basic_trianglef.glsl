#version 330

layout (location = 0) out vec4 color;
in vec3 out_color;
uniform float transperancy;

void main() {
	color = vec4(out_color, transperancy);
}