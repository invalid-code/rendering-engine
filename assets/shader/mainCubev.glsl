#version 330 core

layout (location = 0) in vec3 pos;
layout (location = 1) in vec3 norm;

out vec3 outNorm;
out vec3 outFragPos;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

void main() {
	gl_Position = projection * view * model * vec4(pos, 1.0);
	outFragPos = vec3(model * vec4(pos, 1.0));
	outNorm = norm;
}
