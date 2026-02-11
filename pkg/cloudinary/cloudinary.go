package cloudinary

import (
	"context"
	"errors"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/iamarpitzala/aca-reca-backend/config"
)

var ErrCloudinaryNotConfigured = errors.New("cloudinary is not configured (missing env: CLOUDINARY_CLOUD_NAME, API_KEY, API_SECRET)")

const (
	ResourceTypeImage        = "image"
	ResourceTypeRaw          = "raw"
	DefaultImageFolder       = "acareca/images"
	DefaultDocumentFolder    = "acareca/documents"
	MaxImageSizeBytes        = 10 << 20   // 10 MB
	MaxDocumentSizeBytes     = 20 << 20   // 20 MB
)

// Service handles uploads to Cloudinary.
type Service struct {
	cld *cloudinary.Cloudinary
}

// UploadResult holds the URL and public ID returned from Cloudinary.
type UploadResult struct {
	URL      string
	PublicID string
}

// NewService creates a Cloudinary service from config. Returns nil if config is incomplete.
func NewService(cfg config.CloudinaryConfig) (*Service, error) {
	if cfg.CloudName == "" || cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, nil
	}
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, err
	}
	return &Service{cld: cld}, nil
}

// Upload sends a file to Cloudinary. resourceType should be ResourceTypeImage or ResourceTypeRaw.
// folder is optional; empty string uses no folder.
func (s *Service) Upload(ctx context.Context, file io.Reader, filename string, resourceType string, folder string) (*UploadResult, error) {
	if s == nil || s.cld == nil {
		return nil, ErrCloudinaryNotConfigured
	}
	params := uploader.UploadParams{
		ResourceType: resourceType,
	}
	if folder != "" {
		params.Folder = folder
	}
	if filename != "" {
		params.PublicID = filename
	}
	result, err := s.cld.Upload.Upload(ctx, file, params)
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

// UploadImage uploads an image to the default images folder.
func (s *Service) UploadImage(ctx context.Context, file io.Reader, filename string) (*UploadResult, error) {
	return s.Upload(ctx, file, filename, ResourceTypeImage, DefaultImageFolder)
}

// UploadDocument uploads a raw file (e.g. PDF) to the default documents folder.
func (s *Service) UploadDocument(ctx context.Context, file io.Reader, filename string) (*UploadResult, error) {
	return s.Upload(ctx, file, filename, ResourceTypeRaw, DefaultDocumentFolder)
}
