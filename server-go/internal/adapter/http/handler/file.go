package handler

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/usecase/apperror"
)

const maxImageSize = 5 << 20 // 5MB

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

type File struct {
	staticPath string
	logger     *slog.Logger
}

func NewFile(staticPath string, logger *slog.Logger) *File {
	return &File{staticPath: staticPath, logger: logger}
}

func (h *File) UploadImage(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return apperror.New(apperror.CodeInvalidParam, "缺少文件字段 file")
	}
	if fh.Size > maxImageSize {
		return apperror.New(apperror.CodeInvalidParam, "图片不能超过 5MB")
	}
	ext := strings.ToLower(path.Ext(fh.Filename))
	if !allowedImageExts[ext] {
		return apperror.New(apperror.CodeInvalidParam, "仅支持 jpg/png/gif/webp")
	}
	relDir := filepath.Join("uploads", time.Now().Format("200601"))
	absDir := filepath.Join(h.staticPath, relDir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上传失败", err)
	}
	name := uuid.NewString() + ext
	dst := filepath.Join(absDir, name)

	src, err := fh.Open()
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上传失败", err)
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上传失败", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上传失败", err)
	}
	url := fmt.Sprintf("/%s/%s", h.staticPath, filepath.ToSlash(filepath.Join(relDir, name)))
	return dto.OK(c, map[string]string{"url": url})
}
