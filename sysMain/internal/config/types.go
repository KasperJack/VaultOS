package config

type SoftwareConfig map[string]SoftwareDetails

type SoftwareDetails struct {
	Portable   bool     `yaml:"portable"`
	Category   string   `yaml:"category"`
	Executable string   `yaml:"executable"`
	Junctions  []string `yaml:"junctions,omitempty"`
}
