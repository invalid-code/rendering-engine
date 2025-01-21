#version 330 core

out vec4 color;

uniform vec3 lightColor;
uniform vec3 objColor;

void main() {
    color = vec4(lightColor * objColor, 1.0);
}
