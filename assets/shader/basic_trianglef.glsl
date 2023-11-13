#version 330

layout (location = 0) out vec4 color;
in vec3 out_color;
uniform float green_color;

void main() {
	color = vec4(out_color.x, green_color, out_color.z, 1.0);
}