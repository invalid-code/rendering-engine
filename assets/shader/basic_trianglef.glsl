#version 330

layout (location = 0) out vec4 color;
in vec3 out_color;
in vec2 out_tex_coord;
uniform sampler2D texture1;
uniform sampler2D texture2;

void main() {
	color = mix(texture(texture1, out_tex_coord), texture(texture2, out_tex_coord), 0.2);
}