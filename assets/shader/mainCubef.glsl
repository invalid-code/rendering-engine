#version 330 core

struct Material {
    sampler2D diffuse;
    sampler2D specular;
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

uniform vec3 viewPos;
uniform sampler2D mainTexture;
uniform Material material;
uniform Light light;

void main() {
    vec3 ambient = light.ambient * vec3(texture(material.diffuse, outTexCoord));

    vec3 norm = normalize(outNorm);
    vec3 lightDir = normalize(light.pos - outFragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, outTexCoord));

    vec3 viewDir = normalize(viewPos - outFragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * spec * vec3(texture(material.specular, outTexCoord));

    vec3 res = ambient + diffuse + specular;
    color = vec4(res, 1.0);
}
