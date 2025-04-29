package config

type Registers struct {
	Sources  []string `mapstructure:"sources" json:"sources" yaml:"sources"`
	Replicas []string `mapstructure:"replicas" json:"replicas" yaml:"replicas"`
	Policy   string   `mapstructure:"policy" json:"policy" yaml:"policy"`
	Tables   []string `mapstructure:"tables" json:"tables" yaml:"tables"`
}
