#version 330 core

struct Material {
    sampler2D diffuse;
    sampler2D specular;
    float shinniness;
};

struct DirLight {
    vec3 direction;
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

struct PointLight {
    vec3 pos;
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float constant;
    float linear;
    float quadratic;
};

struct FlashLight {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    vec3 pos;
    vec3 direction;
    float cutOff;
    float outerCutOff;
    float constant;
    float linear;
    float quadratic;
};

out vec4 color;

in vec3 outNorm;
in vec3 outFragPos;
in vec2 outTexCoord;

uniform vec3 viewPos;
uniform sampler2D mainTexture;
uniform Material material;
uniform FlashLight light;

void main() {
    vec3 lightDir = normalize(light.pos - outFragPos);
    float theta = dot(lightDir, normalize(-light.direction));
    float epsilon = light.cutOff - light.outerCutOff;
    float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);
    
    vec3 ambient = light.ambient * vec3(texture(material.diffuse, outTexCoord));

    vec3 norm = normalize(outNorm);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, outTexCoord));

    vec3 viewDir = normalize(viewPos - outFragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * spec * vec3(texture(material.specular, outTexCoord));

    float distance = length(light.pos - outFragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    diffuse *= intensity;
    specular *= intensity;
    diffuse *= attenuation;
    specular *= attenuation;

    vec3 res = ambient + diffuse + specular;
    color = vec4(res, 1.0);
}
