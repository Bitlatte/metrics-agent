package types

type MetricType string

const (
	MetricCPUModel         MetricType = "cpu.model"
	MetricCPUUsage         MetricType = "cpu.usage"
	MetricCPUCountPhysical MetricType = "cpu.count_physical"
	MetricCPUCountLogical  MetricType = "cpu.count_logical"

	MetricMemoryUsage MetricType = "memory.usage"
	MetricMemoryTotal MetricType = "memory.total"
	MetricMemoryFree  MetricType = "memory.free"
)

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeFloat  ValueType = "float"
	ValueTypeBool   ValueType = "bool"
	ValueTypeInt    ValueType = "int"
)

type Labels map[string]string

type Metric struct {
	Type        MetricType `json:"type"`
	ValueType   ValueType  `json:"value_type"`
	StringValue *string    `json:"string_value"`
	FloatValue  *float64   `json:"float_value"`
	IntValue    *int       `json:"int_value"`
	BoolValue   *bool      `json:"bool_value"`
	Unit        string     `json:"unit,omitempty"`
	Timestamp   int64      `json:"timestamp"`
	Labels      Labels     `json:"labels,omitempty"`
}
