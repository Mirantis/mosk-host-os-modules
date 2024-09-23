package domain

type (
	// HostOSConfigurationModules emulates the CRD HostOSConfigurationModules
	// object from the kaas/core.
	HostOSConfigurationModules struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Spec struct {
			Modules []Module `yaml:"modules"`
		} `yaml:"spec"`
	}

	// Module is a minimal required structure to represent a module.
	Module struct {
		NameVersionTuple `yaml:",inline"`
		Sha256Sum        string `yaml:"sha256sum"`
	}

	// NameVersionTuple represents a pair of name-version both required
	// during deserialization from YAML format.
	NameVersionTuple struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
)

func (m HostOSConfigurationModules) IsEmpty() bool {
	return m.APIVersion == "" || m.Kind == "" || m.Metadata.Name == "" ||
		len(m.Spec.Modules) == 0
}

func (t NameVersionTuple) String() string {
	return str(t.Name, t.Version)
}

func (m Module) IsEqual(a Module) bool {
	return m.NameVersionTuple == a.NameVersionTuple && m.Sha256Sum == a.Sha256Sum
}

func (m Module) String() string {
	return str(m.Name, m.Version)
}

func str(n, v string) string {
	return n + "-" + v
}
