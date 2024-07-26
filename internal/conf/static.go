package conf

import (
	"net/url"
	"os"
)

// CustomConf returns the absolute path of custom configuration file that is used.
var CustomConf string

// Build time and commit information.
//
// ⚠️ WARNING: should only be set by "-ldflags".
var (
	BuildTime   string
	BuildCommit string
)

var (
	App struct {
		// ⚠️ WARNING: Should only be set by the main package (i.e. "imail.go").
		Version string `ini:"-"`

		Name      string
		BrandName string
		RunUser   string
		RunMode   string
	}

	// log
	Log struct {
		Format   string
		RootPath string
	}

	// log
	Admin struct {
		User   string
		Pass string
	}

	Plugins struct{
		PathDir string
	}

	// Cache settings
	Cache struct {
		Adapter  string
		Interval int
		Host     string
	}

	// web settings
	Web struct {
		HttpAddr                 string `ini:"http_addr"`
		HttpPort                 int    `ini:"http_port"`
		Domain                   string
		AppDataPath              string
		AccessControlAllowOrigin string

		ExternalURL          string `ini:"EXTERNAL_URL"`
		Protocol             string
		CertFile             string
		KeyFile              string
		TLSMinVersion        string `ini:"TLS_MIN_VERSION"`
		UnixSocketPermission string
		LocalRootURL         string `ini:"LOCAL_ROOT_URL"`

		OfflineMode      bool
		DisableRouterLog bool
		EnableGzip       bool

		LoadAssetsFromDisk bool

		LandingURL string `ini:"LANDING_URL"`

		// Derived from other static values
		URL            *url.URL    `ini:"-"` // Parsed URL object of ExternalURL.
		Subpath        string      `ini:"-"` // Subpath found the ExternalURL. Should be empty when not found.
		SubpathDepth   int         `ini:"-"` // The number of slashes found in the Subpath.
		UnixSocketMode os.FileMode `ini:"-"` // Parsed file mode of UnixSocketPermission.

		MailSaveMode string
	}

	// Authentication settings
	Auth struct {
		ActivateCodeLives         int
		ResetPasswordCodeLives    int
		RequireEmailConfirmation  bool
		RequireSigninView         bool
		DisableRegistration       bool
		EnableRegistrationCaptcha bool

		EnableReverseProxyAuthentication   bool
		EnableReverseProxyAutoRegistration bool
		ReverseProxyAuthenticationHeader   string
	}

	// Session settings
	Session struct {
		Provider       string
		ProviderConfig string
		CookieName     string
		CookieSecure   bool
		GCInterval     int64 `ini:"gc_interval"`
		MaxLifeTime    int64
		CSRFCookieName string `ini:"csrf_cookie_name"`
	}


	// Security settings
	Security struct {
		InstallLock             bool
		SecretKey               string
		LoginRememberDays       int
		CookieRememberName      string
		CookieUsername          string
		CookieSecure            bool
		EnableLoginStatusCookie bool
		LoginStatusCookieName   string
	}

)
