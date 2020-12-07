package configs

// HashConfig type
type HashConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultHashConfig function to generate default hash config
func DefaultHashConfig() HashConfig {
	return HashConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
}
