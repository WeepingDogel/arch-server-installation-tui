package mirror

// Mirror represents a pacman mirror server.
type Mirror struct {
	Name    string
	URL     string
	Country string
	Special bool // Whether this is a special mirror (Arch Linux CN, etc.)
}

// DefaultMirrors returns the complete list of Arch Linux mirrors.
// Includes special mirrors like Arch Linux CN, TUNA, USTC, 163, etc.
func DefaultMirrors() []Mirror {
	return []Mirror{
		// Global official mirrors
		{Name: "Arch Linux Official (Global)", URL: "http://mirror.archlinuxarm.org", Country: "Global", Special: false},
		{Name: "United States - Kernel.org", URL: "http://mirrors.kernel.org/archlinux", Country: "US", Special: false},
		{Name: "United States - Princeton", URL: "http://mirror.math.princeton.edu/pub/archlinux", Country: "US", Special: false},
		{Name: "Germany - GWDG", URL: "https://archlinux.gwdg.de", Country: "DE", Special: false},

		// China - Special mirrors (Arch Linux CN)
		{Name: "China - TUNA (Tsinghua) ★", URL: "https://mirrors.tuna.tsinghua.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - USTC (USTC) ★", URL: "https://mirrors.ustc.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - 163 (NetEase) ★", URL: "https://mirrors.163.com/archlinux", Country: "CN", Special: true},
		{Name: "China - Aliyun (Alibaba Cloud) ★", URL: "https://mirrors.aliyun.com/archlinux", Country: "CN", Special: true},
		{Name: "China - Huawei Cloud ★", URL: "https://mirrors.huaweicloud.com/archlinux", Country: "CN", Special: true},
		{Name: "China - Tencent Cloud ★", URL: "https://mirrors.tencent.com/archlinux", Country: "CN", Special: true},
		{Name: "China - SJTU (Shanghai Jiao Tong) ★", URL: "https://mirrors.sjtug.sjtu.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - Nanjing University ★", URL: "https://mirrors.nju.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - Chongqing University ★", URL: "https://mirrors.cqu.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - Beijing Foreign Studies ★", URL: "https://mirrors.bfsu.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - Neusoft ★", URL: "https://mirrors.neusoft.edu.cn/archlinux", Country: "CN", Special: true},
		{Name: "China - Xiyou Linux ★", URL: "https://mirrors.xiyoulinux.cn/archlinux", Country: "CN", Special: true},

		// Asia Pacific
		{Name: "Japan - JAIST", URL: "http://ftp.jaist.ac.jp/pub/Linux/ArchLinux", Country: "JP", Special: false},
		{Name: "Singapore - 0x", URL: "https://mirror.0x.sg/archlinux", Country: "SG", Special: false},
		{Name: "Taiwan - NCHC", URL: "https://ftp.archlinux.tw/archlinux", Country: "TW", Special: false},
		{Name: "South Korea - KAIST", URL: "http://mirror.kaist.edu.cn/archlinux", Country: "KR", Special: false},
		{Name: "India - IIT Bombay", URL: "http://mirror.cse.iitb.ac.in/archlinux", Country: "IN", Special: false},

		// Europe
		{Name: "Netherlands - NLUUG", URL: "https://ftp.nluug.nl/pub/os/Linux/distr/archlinux", Country: "NL", Special: false},
		{Name: "France - IRCAM", URL: "https://mirrors.ircam.fr/pub/archlinux", Country: "FR", Special: false},
		{Name: "UK - UKFast", URL: "https://mirror.ukfast.co.uk/sites/archlinux.org", Country: "GB", Special: false},
		{Name: "Sweden - Lysator", URL: "https://ftp.lysator.liu.se/pub/archlinux", Country: "SE", Special: false},
		{Name: "Poland - ICM", URL: "https://mirror.icm.edu.pl/archlinux", Country: "PL", Special: false},
		{Name: "Russia - Yandex", URL: "https://mirror.yandex.ru/archlinux", Country: "RU", Special: false},
		{Name: "Austria - VUM", URL: "https://mirror.vunet.eu/archlinux", Country: "AT", Special: false},

		// Arch Linux CN repository support (for AUR/Chinese packages)
		{Name: "Arch Linux CN - TUNA ★", URL: "https://mirrors.tuna.tsinghua.edu.cn/archlinuxcn", Country: "CN", Special: true},
		{Name: "Arch Linux CN - USTC ★", URL: "https://mirrors.ustc.edu.cn/archlinuxcn", Country: "CN", Special: true},
		{Name: "Arch Linux CN - 163 ★", URL: "https://mirrors.163.com/archlinux-cn", Country: "CN", Special: true},
	}
}

// FilterByCountry returns mirrors matching the given country code.
func FilterByCountry(mirrors []Mirror, country string) []Mirror {
	if country == "" || country == "all" {
		return mirrors
	}
	var filtered []Mirror
	for _, m := range mirrors {
		if m.Country == country {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// FilterSpecial returns only special mirrors (Arch Linux CN, etc.).
func FilterSpecial(mirrors []Mirror) []Mirror {
	var filtered []Mirror
	for _, m := range mirrors {
		if m.Special {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// SearchByName searches mirrors by name substring.
func SearchByName(mirrors []Mirror, query string) []Mirror {
	if query == "" {
		return mirrors
	}
	var filtered []Mirror
	for _, m := range mirrors {
		if containsIgnoreCase(m.Name, query) || containsIgnoreCase(m.URL, query) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// containsIgnoreCase checks if substr is in s (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	substr = toLower(substr)
	s = toLower(s)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// toLower converts a string to lowercase (avoiding unicode import).
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + 32
		}
		result[i] = c
	}
	return string(result)
}