package config

type Config struct {
    ProxyList          []string
    UserAgents         []string
    TimeoutSeconds     int
    MaxRetries         int
    ConcurrentRequests int
}

func GetDefaultConfig() *Config {
    return &Config{
        TimeoutSeconds:     10,
        MaxRetries:         3,
        ConcurrentRequests: 20,
        UserAgents: []string{
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/121.0.0.0 Safari/537.36",
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/121.0.0.0 Safari/537.36",
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
        },
    }
}
