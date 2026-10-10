package configs

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"aarambha-agency/internal/domain"
)

type Config struct {
	Host              string
	Port              string
	Env               string
	DatabaseURL       string
	SupabaseURL       string
	SupabaseSecretKey string
	SupabaseBucket    string
	SessionSecret     string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Company           domain.CompanyInfo
}

func LoadConfig() *Config {
	loadDotEnv(".env")

	host := getEnv("HOST", "0.0.0.0")
	port := getEnv("PORT", "5196")
	env := getEnv("APP_ENV", "development")

	dbURL := getEnv("DB_URL", "")
	supabaseURL := getEnv("SUPABASE_URL", "")
	supabaseSecretKey := getEnv("SUPABASE_SECRET_KEY", "")
	supabaseBucket := getEnv("SUPABASE_BUCKET", "test-bkt")
	sessionSecret := getEnv("SESSION_SECRET", getEnv("SUPABASE_SECRET_KEY", "aarambha-agency-auth-salt-7386985852"))

	readTimeoutSec, _ := strconv.Atoi(getEnv("READ_TIMEOUT_SEC", "10"))
	writeTimeoutSec, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT_SEC", "15"))
	idleTimeoutSec, _ := strconv.Atoi(getEnv("IDLE_TIMEOUT_SEC", "60"))

	return &Config{
		Host:              host,
		Port:              port,
		Env:               env,
		DatabaseURL:       dbURL,
		SupabaseURL:       supabaseURL,
		SupabaseSecretKey: supabaseSecretKey,
		SupabaseBucket:    supabaseBucket,
		SessionSecret:     sessionSecret,
		ReadTimeout:       time.Duration(readTimeoutSec) * time.Second,
		WriteTimeout:      time.Duration(writeTimeoutSec) * time.Second,
		IdleTimeout:       time.Duration(idleTimeoutSec) * time.Second,
		Company: domain.CompanyInfo{
			Name:      "Aarambha",
			BrandName: "AARAMBHA <span>AGENCY</span>",
			Tagline:   "Curated Indian Real Estate & Ethical Advisory",
			Headquarters: domain.Address{
				Address: "Level 4, Skyview Corporate Park, Hitec City",
				City:    "Hyderabad",
				State:   "Telangana",
				Pincode: "500081",
				Country: "India",
			},
			Branches: []domain.Branch{
				{City: "Hyderabad", Area: "Hitec City & Jubilee Hills", Phone: "+91 73869 85852"},
				{City: "Bengaluru", Area: "Indiranagar & Whitefield", Phone: "+91 80 4500 1234"},
				{City: "Mumbai", Area: "Bandra Kurla Complex (BKC)", Phone: "+91 22 6700 9876"},
				{City: "Delhi NCR", Area: "Golf Course Road, Gurugram", Phone: "+91 124 4900 567"},
			},
			Phone:            "+91 73869 85852",
			PhoneRaw:         "+917386985852",
			WhatsApp:         "917386985852",
			Email:            "advisory@aarambha.in",
			Founded:          "Started 1 Year Ago (2025)",
			EstablishedYear:  2025,
			ReraBrokerNumber: "RERA Reg. No: A52000012345 (Telangana State Real Estate Regulatory Authority)",
			Hours:            "Mon - Sat: 9:30 AM to 7:30 PM IST (Sunday by Appointment)",
		},
	}
}

func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// loadDotEnv reads key=value pairs from a file into os.Environ if not already set.
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return // Ignore if .env is missing
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// Strip surrounding quotes if any
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}
}
