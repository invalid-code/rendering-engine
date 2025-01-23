#version 330 core

struct Material {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float shinniness;
};

struct Light {
    vec3 pos;
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

out vec4 color;

in vec3 outNorm;
in vec3 outFragPos;
in vec2 outTexCoord;

uniform vec3 lightColor;
uniform vec3 objColor;
uniform vec3 viewPos;
uniform sampler2D mainTexture;
uniform Material material;
uniform Light light;

void main() {
    float ambientStrength = 0.1;
    vec3 ambient = light.ambient * material.ambient;

    vec3 norm = normalize(outNorm);
    vec3 lightDir = normalize(light.pos - outFragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * (diff * material.diffuse);

    float specStrength = 0.5;
    vec3 viewDir = normalize(viewPos - outFragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * (spec * material.specular);

    vec3 res = (ambient + diffuse + specular) * objColor;
    color = texture(mainTexture, outTexCoord);
}
