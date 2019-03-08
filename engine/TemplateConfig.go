package engine

type TemplateConfig struct {
	LayoutPath       string `json:"layout_path"`
	IncludePath      string `json:"include_path"`
	IncludeCondition string `json:"include_condition"`
}
