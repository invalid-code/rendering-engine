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

struct SpotLight {
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
uniform DirLight dirLight;
#define NO_POINT_LIGHTS 4
uniform PointLight pointLights[NO_POINT_LIGHTS];
uniform SpotLight spotLight;

vec3 calcDirLight(DirLight light, vec3 norm, vec3 viewDir) {
    vec3 lightDir = normalize(-light.direction);

    vec3 ambient = light.ambient * vec3(texture(material.diffuse, outTexCoord));
    
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, outTexCoord));

    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * spec * vec3(texture(material.specular, outTexCoord));
    
    return ambient + diffuse + specular;
}

vec3 calcPointLight(PointLight light, vec3 norm, vec3 fragPos, vec3 viewDir) {
    vec3 lightDir = normalize(light.pos - outFragPos);

    vec3 ambient = light.ambient * vec3(texture(material.diffuse, outTexCoord));
    
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, outTexCoord));

    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * spec * vec3(texture(material.specular, outTexCoord));

    float distance = length(light.pos - outFragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;
    return ambient + diffuse + specular;
}

vec3 calcSpotLight(SpotLight light, vec3 norm, vec3 fragPos, vec3 viewDir) {
    vec3 lightDir = normalize(light.pos - outFragPos);

    vec3 ambient = light.ambient * vec3(texture(material.diffuse, outTexCoord));
    
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, outTexCoord));

    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shinniness);
    vec3 specular = light.specular * spec * vec3(texture(material.specular, outTexCoord));

    float theta = dot(lightDir, normalize(-light.direction));
    float epsilon = light.cutOff - light.outerCutOff;
    float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);

    float distance = length(light.pos - outFragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    ambient *= intensity * attenuation;
    diffuse *= intensity * attenuation;
    specular *= intensity * attenuation;
    return ambient + diffuse + specular;
}

void main() {
    vec3 norm = normalize(outNorm);
    vec3 viewDir = normalize(viewPos - outFragPos);

    vec3 res = calcDirLight(dirLight, norm, viewDir);
    for (int i = 0; i < NO_POINT_LIGHTS; i++)
        res += calcPointLight(pointLights[i], outNorm, outFragPos, viewDir);
    res += calcSpotLight(spotLight, norm, outFragPos, viewDir);

    color = vec4(res, 1.0);
}
