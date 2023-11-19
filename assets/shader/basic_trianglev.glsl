#version 330

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 color;
layout (location = 2) in vec2 tex_coord;
out vec3 out_color;
out vec2 out_tex_coord;
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

void main() {	
	gl_Position = projection * view * model * vec4(position, 1.0);
	out_color = color;
	out_tex_coord = tex_coord;
}
