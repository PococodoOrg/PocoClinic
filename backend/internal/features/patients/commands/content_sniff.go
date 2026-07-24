package commands

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
)

func sniffAndValidateContent(content io.Reader, declaredType, fileName string) (string, io.Reader, error) {
	prefix := make([]byte, 512)
	n, err := io.ReadFull(content, prefix)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "Failed to read uploaded file")
	}
	prefix = prefix[:n]
	if n == 0 {
		return "", nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File is empty")
	}

	detected := strings.TrimSpace(strings.Split(http.DetectContentType(prefix), ";")[0])
	if _, ok := domain.AllowedDocumentContentTypes[detected]; !ok {
		return "", nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File type not allowed. Use PDF, PNG, JPEG, GIF, or plain text.")
	}

	declared := normalizeContentType(declaredType, fileName)
	if declared != "" && declared != "application/octet-stream" && declared != detected {
		return "", nil, pkgerrors.NewAPIError(pkgerrors.ErrValidation, "File content does not match declared type")
	}

	return detected, io.MultiReader(bytes.NewReader(prefix), content), nil
}
