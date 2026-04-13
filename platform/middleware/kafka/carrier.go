package kafka

// kafkaHeaderCarrier адаптирует map[string][]byte под otel TextMapCarrier
type kafkaHeaderCarrier map[string][]byte

func (c kafkaHeaderCarrier) Get(key string) string {
	if val, ok := c[key]; ok {
		return string(val)
	}
	return ""
}

func (c kafkaHeaderCarrier) Set(key, val string) {
	c[key] = []byte(val)
}

func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}