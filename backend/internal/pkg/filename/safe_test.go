package filename_test

import (
	"strings"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/filename"
)

func TestSafeAttachmentFilename(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"report.pdf", "report.pdf"},
		{`../../etc/passwd`, "passwd"},
		{`bad"; evil`, "bad evil"},
		{"", "download"},
	}

	for _, tt := range tests {
		got := filename.SafeAttachmentFilename(tt.in)
		if got != tt.want {
			t.Fatalf("SafeAttachmentFilename(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if strings.ContainsAny(got, "\";\\") {
			t.Fatalf("unsafe characters in %q", got)
		}
	}
}
