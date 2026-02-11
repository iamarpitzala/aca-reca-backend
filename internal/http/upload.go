package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/pkg/cloudinary"
)

// Allowed image and document types for upload validation.
var (
	allowedImageTypes = map[string]bool{
		"image/jpeg": true, "image/jpg": true, "image/png": true,
		"image/gif": true, "image/webp": true,
	}
	allowedDocumentTypes = map[string]bool{
		"application/pdf": true,
		"application/msword": true, "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
)

// UploadHandler handles image and document uploads to Cloudinary.
type UploadHandler struct {
	svc *cloudinary.Service
}

// NewUploadHandler creates an upload handler. svc may be nil if Cloudinary is not configured.
func NewUploadHandler(svc *cloudinary.Service) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// UploadImage godoc
// @Summary Upload an image
// @Description Upload a single image file to Cloudinary (max 10 MB). Allowed: JPEG, PNG, GIF, WebP.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} map[string]string "url, public_id"
// @Failure 400 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	if h.svc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "upload service is not configured"})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid file: " + err.Error()})
		return
	}
	defer file.Close()

	if header.Size > cloudinary.MaxImageSizeBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "image must be at most 10 MB"})
		return
	}
	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image type; allowed: JPEG, PNG, GIF, WebP"})
		return
	}

	filename := safePublicID(header.Filename, "img")
	result, err := h.svc.UploadImage(c.Request.Context(), file, filename)
	if err != nil {
		if err == cloudinary.ErrCloudinaryNotConfigured {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": result.URL, "public_id": result.PublicID})
}

// UploadDocument godoc
// @Summary Upload a document
// @Description Upload a single document to Cloudinary (max 20 MB). Allowed: PDF, DOC, DOCX.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Success 200 {object} map[string]string "url, public_id"
// @Failure 400 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /upload/document [post]
func (h *UploadHandler) UploadDocument(c *gin.Context) {
	if h.svc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "upload service is not configured"})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid file: " + err.Error()})
		return
	}
	defer file.Close()

	if header.Size > cloudinary.MaxDocumentSizeBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "document must be at most 20 MB"})
		return
	}
	contentType := header.Header.Get("Content-Type")
	if !allowedDocumentTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type; allowed: PDF, DOC, DOCX"})
		return
	}

	filename := safePublicID(header.Filename, "doc")
	result, err := h.svc.UploadDocument(c.Request.Context(), file, filename)
	if err != nil {
		if err == cloudinary.ErrCloudinaryNotConfigured {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": result.URL, "public_id": result.PublicID})
}

// safePublicID returns a safe public_id: optional prefix + UUID + extension (if any).
func safePublicID(originalFilename, prefix string) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext != "" {
		ext = strings.TrimPrefix(ext, ".")
		// Cloudinary accepts alphanumeric and _ for public_id; keep extension short.
		if len(ext) > 4 {
			ext = ext[:4]
		}
		return prefix + "_" + uuid.New().String() + "." + ext
	}
	return prefix + "_" + uuid.New().String()
}
