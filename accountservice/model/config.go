package model

type SpringCloudConfig struct {
        Name string `json:"name"`
        Profiles []string `json:"profiles"`
        Label string `json:"label"`
        Version string `json:"version"`
        PropertySources []PropertySource `json:"propertySources"`
}

type PropertySource struct {
        Name string `json:"name"`
        Source map[string]interface{} `json:"source"`
}
