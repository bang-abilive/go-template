package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// allConfigVars là danh sách tất cả các biến môi trường được đọc bởi Config.
// Dùng để đảm bảo cách ly giữa các test khi godotenv gọi os.Setenv trực tiếp.
var allConfigVars = []string{
	"APP_ENV", "APP_DEBUG", "APP_NAME", "APP_VERSION",
	"SERVER_PORT",
	"DB_NAME", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
	"DB_SSL_MODE", "DB_MAX_CONNS", "DB_MAX_IDLE_CONNS", "DB_MAX_LIFETIME_CONNS",
}

// isolateEnv lưu trạng thái hiện tại của tất cả config vars và khôi phục sau test.
// Cần thiết vì godotenv gọi os.Setenv trực tiếp, bypass cơ chế cleanup của t.Setenv.
func isolateEnv(t *testing.T) {
	t.Helper()
	for _, v := range allConfigVars {
		old, exists := os.LookupEnv(v)
		_ = os.Unsetenv(v)
		v := v
		t.Cleanup(func() {
			if exists {
				_ = os.Setenv(v, old)
			} else {
				_ = os.Unsetenv(v)
			}
		})
	}
}

// setValidEnv đặt tất cả các biến môi trường bắt buộc để tạo ra một Config hợp lệ.
// DB_USER và DB_PASSWORD có tag "unset" nên sẽ bị xóa khỏi env sau khi parse;
// không nên assert giá trị của chúng sau khi gọi LoadFromEnv.
func setValidEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpassword")
}

// TestServerAddress kiểm tra định dạng địa chỉ server trả về.
func TestServerAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		port     uint16
		expected string
	}{
		{name: "default port 8080", port: 8080, expected: ":8080"},
		{name: "custom port 3000", port: 3000, expected: ":3000"},
		{name: "port 443 HTTPS", port: 443, expected: ":443"},
		{name: "port 0", port: 0, expected: ":0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := Config{Server: ServerConfig{Port: tt.port}}
			assert.Equal(t, tt.expected, cfg.Server.ServerAddress())
		})
	}
}

// TestLoadFromEnv_DefaultValues kiểm tra rằng các giá trị mặc định được áp dụng đúng
// khi không đặt các biến tùy chọn.
func TestLoadFromEnv_DefaultValues(t *testing.T) {
	t.Chdir(t.TempDir())
	isolateEnv(t)
	setValidEnv(t)

	cfg, err := LoadFromEnv()
	must := require.New(t)
	is := assert.New(t)

	must.NoError(err)
	must.NotNil(cfg)

	is.Equal("development", cfg.App.Env)
	is.False(cfg.App.Debug)
	is.Equal("0.0.1", cfg.App.Version)
	is.Empty(cfg.App.Name)
	is.Equal("disable", cfg.Database.SSLMode)
	is.Equal(10, cfg.Database.MaxConns)
	is.Equal(5, cfg.Database.MaxIdleConns)
	is.Equal(time.Hour, cfg.Database.MaxLifetimeConns)
	is.Equal(uint16(8080), cfg.Server.Port)
	is.Equal("testdb", cfg.Database.Name)
	is.Equal("localhost", cfg.Database.Host)
	is.Equal("5432", cfg.Database.Port)
}

// TestLoadFromEnv_CustomValues kiểm tra rằng các giá trị tùy chỉnh được parse đúng.
func TestLoadFromEnv_CustomValues(t *testing.T) {
	t.Chdir(t.TempDir())
	isolateEnv(t)
	setValidEnv(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_DEBUG", "true")
	t.Setenv("APP_NAME", "my-app")
	t.Setenv("APP_VERSION", "1.2.3")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DB_SSL_MODE", "require")
	t.Setenv("DB_MAX_CONNS", "20")
	t.Setenv("DB_MAX_IDLE_CONNS", "3")
	t.Setenv("DB_MAX_LIFETIME_CONNS", "30m")

	cfg, err := LoadFromEnv()
	must := require.New(t)
	is := assert.New(t)

	must.NoError(err)
	must.NotNil(cfg)

	is.Equal("production", cfg.App.Env)
	is.True(cfg.App.Debug)
	is.Equal("my-app", cfg.App.Name)
	is.Equal("1.2.3", cfg.App.Version)
	is.Equal(uint16(9090), cfg.Server.Port)
	is.Equal("require", cfg.Database.SSLMode)
	is.Equal(20, cfg.Database.MaxConns)
	is.Equal(3, cfg.Database.MaxIdleConns)
	is.Equal(30*time.Minute, cfg.Database.MaxLifetimeConns)
}

// TestLoadFromEnv_InvalidAppEnv kiểm tra rằng các giá trị APP_ENV không hợp lệ
// (chứa ký tự đặc biệt, dấu gạch chéo, khoảng trắng) bị từ chối.
func TestLoadFromEnv_InvalidAppEnv(t *testing.T) {
	tests := []struct {
		name   string
		appEnv string
	}{
		{name: "path traversal", appEnv: "../etc/passwd"},
		{name: "forward slash", appEnv: "prod/us"},
		{name: "space", appEnv: "prod env"},
		{name: "special chars", appEnv: "prod!@#"},
		{name: "dot separator", appEnv: "prod.local"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			isolateEnv(t)
			t.Setenv("APP_ENV", tt.appEnv)

			cfg, err := LoadFromEnv()
			is := assert.New(t)

			is.Error(err)
			is.Nil(cfg)
			is.ErrorContains(err, "invalid APP_ENV value")
		})
	}
}

// TestLoadFromEnv_ValidAppEnvNames kiểm tra rằng các tên APP_ENV hợp lệ
// (chữ cái, số, dấu gạch ngang, gạch dưới) được chấp nhận.
func TestLoadFromEnv_ValidAppEnvNames(t *testing.T) {
	tests := []struct {
		name   string
		appEnv string
	}{
		{name: "lowercase", appEnv: "development"},
		{name: "with hyphen", appEnv: "prod-us"},
		{name: "with underscore", appEnv: "prod_us"},
		{name: "uppercase", appEnv: "PRODUCTION"},
		{name: "mixed case with digits", appEnv: "staging1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			isolateEnv(t)
			t.Setenv("APP_ENV", tt.appEnv)
			setValidEnv(t)

			cfg, err := LoadFromEnv()
			must := require.New(t)

			must.NoError(err)
			must.NotNil(cfg)
			assert.Equal(t, tt.appEnv, cfg.App.Env)
		})
	}
}

// TestLoadFromEnv_MissingRequiredFields kiểm tra rằng thiếu các trường bắt buộc
// trong cấu hình database sẽ trả về lỗi.
func TestLoadFromEnv_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T)
	}{
		{
			name: "empty DB_NAME",
			setup: func(t *testing.T) {
				t.Setenv("DB_NAME", "")
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("DB_PORT", "5432")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "pass")
			},
		},
		{
			name: "empty DB_HOST",
			setup: func(t *testing.T) {
				t.Setenv("DB_NAME", "testdb")
				t.Setenv("DB_HOST", "")
				t.Setenv("DB_PORT", "5432")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "pass")
			},
		},
		{
			name: "empty DB_PORT",
			setup: func(t *testing.T) {
				t.Setenv("DB_NAME", "testdb")
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("DB_PORT", "")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "pass")
			},
		},
		{
			name: "empty DB_USER",
			setup: func(t *testing.T) {
				t.Setenv("DB_NAME", "testdb")
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("DB_PORT", "5432")
				t.Setenv("DB_USER", "")
				t.Setenv("DB_PASSWORD", "pass")
			},
		},
		{
			name: "empty DB_PASSWORD",
			setup: func(t *testing.T) {
				t.Setenv("DB_NAME", "testdb")
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("DB_PORT", "5432")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			isolateEnv(t)
			tt.setup(t)

			cfg, err := LoadFromEnv()
			is := assert.New(t)

			is.Error(err)
			is.Nil(cfg)
			is.ErrorContains(err, "[config] parse env")
		})
	}
}

// TestLoadFromEnv_LoadsDotEnvFile kiểm tra rằng các biến từ file .env
// được nạp vào cấu hình khi không có biến môi trường nào được đặt trước.
func TestLoadFromEnv_LoadsDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	isolateEnv(t)

	envContent := "DB_NAME=filedb\nDB_HOST=filehost\nDB_PORT=5433\nDB_USER=fileuser\nDB_PASSWORD=filepass\n"
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(envContent), 0o600)
	require.NoError(t, err)

	cfg, err := LoadFromEnv()
	must := require.New(t)
	is := assert.New(t)

	must.NoError(err)
	must.NotNil(cfg)

	is.Equal("filedb", cfg.Database.Name)
	is.Equal("filehost", cfg.Database.Host)
	is.Equal("5433", cfg.Database.Port)
}

// TestLoadFromEnv_UnreadableEnvFile kiểm tra nhánh cảnh báo khi godotenv.Load thất bại.
// Tạo một thư mục (directory) có tên ".env" — os.Stat thành công nhưng
// io.Copy trả về EISDIR khi godotenv cố đọc, kích hoạt slog.Warn + continue.
// Config vẫn parse thành công từ các env var đã được set trước.
func TestLoadFromEnv_UnreadableEnvFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	isolateEnv(t)
	setValidEnv(t)

	require.NoError(t, os.Mkdir(filepath.Join(dir, ".env"), 0o755))

	cfg, err := LoadFromEnv()
	must := require.New(t)

	must.NoError(err)
	must.NotNil(cfg)
	assert.Equal(t, "testdb", cfg.Database.Name)
}

// TestLoadFromEnv_EnvFileOverrideOrder kiểm tra rằng file .env.{mode}.local
// có độ ưu tiên cao hơn .env vì được load trước.
// godotenv.Load KHÔNG ghi đè biến đã tồn tại, nên file đầu tiên được load sẽ "thắng".
func TestLoadFromEnv_EnvFileOverrideOrder(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	isolateEnv(t)
	t.Setenv("APP_ENV", "test")

	baseEnv := "DB_NAME=basedb\nDB_HOST=basehost\nDB_PORT=5432\nDB_USER=baseuser\nDB_PASSWORD=basepass\nAPP_VERSION=1.0.0\n"
	localEnv := "DB_NAME=localdb\nDB_HOST=localhost\nDB_PORT=5555\nDB_USER=localuser\nDB_PASSWORD=localpass\nAPP_VERSION=2.0.0\n"

	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(baseEnv), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env.test.local"), []byte(localEnv), 0o600))

	cfg, err := LoadFromEnv()
	must := require.New(t)
	is := assert.New(t)

	must.NoError(err)
	must.NotNil(cfg)

	// .env.test.local được load trước nên các giá trị của nó "thắng".
	// Khi .env load sau, DB_NAME đã tồn tại nên godotenv bỏ qua.
	is.Equal("localdb", cfg.Database.Name)
	is.Equal("2.0.0", cfg.App.Version)
}
