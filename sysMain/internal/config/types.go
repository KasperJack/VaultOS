package config

// SoftwareConfig represents the structure of the YAML configuration file
type SoftwareConfig map[string]SoftwareDetails

// SoftwareDetails holds the details of a software package
type SoftwareDetails struct {
	Portable   bool     `yaml:"portable"`
	Category   string   `yaml:"category"`
	Executable string   `yaml:"executable"`
	Junctions  []string `yaml:"junctions,omitempty"`
}
