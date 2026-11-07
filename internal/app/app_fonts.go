package app

import (
	"regexp"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// GetInstalledFonts читает список реально установленных шрифтов Windows из реестра
func (a *App) GetInstalledFonts() ([]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.QUERY_VALUE)
	if err != nil {
		// Фолбек на моноширинные и системные шрифты, если реестр недоступен
		return []string{"Cascadia Code", "Consolas", "Courier New", "Liberation Mono"}, nil
	}
	defer key.Close()

	names, err := key.ReadValueNames(-1)
	if err != nil {
		return []string{"Cascadia Code", "Consolas", "Courier New"}, nil
	}

	trailingTypeRegex := regexp.MustCompile(`\s*\((TrueType|OpenType|PostScript|Type 1)\)\s*$`)
	fontSet := make(map[string]struct{})

	for _, name := range names {
		cleanName := trailingTypeRegex.ReplaceAllString(name, "")
		cleanName = strings.TrimSpace(cleanName)
		if cleanName != "" {
			fontSet[cleanName] = struct{}{}
		}
	}

	var result []string
	for font := range fontSet {
		result = append(result, font)
	}

	sort.Strings(result)
	return result, nil
}
