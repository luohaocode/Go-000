package config

import "testing"

func TestApplyYAML(t *testing.T) {
	c := new(DbConfig)
	if err := ApplyYAML(GetYAML("./demo.yaml"), c); err != nil {
		t.Errorf("Apply Yaml error")
	}

	t.Logf("%v", c)
}
