package filename

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// SafeAttachmentFilename returns a basename safe for Content-Disposition headers.
func SafeAttachmentFilename(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range name {
		if r < 32 || r == '"' || r == ';' || r == '\\' || unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	name = strings.TrimSpace(b.String())
	if name == "" || name == "." {
		return "download"
	}
	return name
}

// ContentDispositionAttachment builds a safe attachment Content-Disposition value.
func ContentDispositionAttachment(name string) string {
	return fmt.Sprintf(`attachment; filename=%q`, SafeAttachmentFilename(name))
}
