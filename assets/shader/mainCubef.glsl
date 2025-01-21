#version 330 core

out vec4 color;

in vec3 outNorm;
in vec3 outFragPos;

uniform vec3 lightColor;
uniform vec3 objColor;
uniform vec3 lightPos;

void main() {
    float ambientStrength = 0.1;
    vec3 ambient = ambientStrength * lightColor;

    vec3 norm = normalize(outNorm);
    vec3 lightDir = normalize(lightPos - outFragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    vec3 res = (ambient + diffuse) * objColor;
    color = vec4(res, 1.0);
}
