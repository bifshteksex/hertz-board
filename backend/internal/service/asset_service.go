package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/repository"
)

const (
	MaxFileSize     = 10 * 1024 * 1024 // 10MB
	ThumbnailWidth  = 300
	ThumbnailHeight = 300
	MaxImageWidth   = 4000
	MaxImageHeight  = 4000
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type AssetService struct {
	assetRepo     *repository.AssetRepository
	workspaceRepo *repository.WorkspaceRepository
	minioClient   *minio.Client
	bucketName    string
	endpoint      string
}

func NewAssetService(
	assetRepo *repository.AssetRepository,
	workspaceRepo *repository.WorkspaceRepository,
	minioEndpoint, minioAccessKey, minioSecretKey string,
	useSSL bool,
) (*AssetService, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	bucketName := "hertz-board-assets"

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy to public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, bucketName)

		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &AssetService{
		assetRepo:     assetRepo,
		workspaceRepo: workspaceRepo,
		minioClient:   minioClient,
		bucketName:    bucketName,
		endpoint:      minioEndpoint,
	}, nil
}

// UploadAsset uploads a file to MinIO and creates an asset record
func (s *AssetService) UploadAsset(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	filename, contentType string,
	size int64,
	reader io.Reader,
) (*models.Asset, error) {
	if err := s.validateUpload(size, contentType); err != nil {
		return nil, err
	}

	fileData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("%s/%s/%s%s", workspaceID, time.Now().Format("2006/01"), uuid.New(), ext)

	isImage := AllowedImageTypes[contentType]
	width, height, thumbnailURL, err := s.processImage(ctx, fileData, contentType, isImage, ext, workspaceID)
	if err != nil {
		return nil, err
	}

	if err := s.uploadFile(ctx, objectName, fileData, size, contentType); err != nil {
		return nil, err
	}

	asset := &models.Asset{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		UploadedBy:   userID,
		Filename:     filename,
		ContentType:  contentType,
		Size:         size,
		URL:          s.getObjectURL(objectName),
		ThumbnailURL: thumbnailURL,
		Width:        width,
		Height:       height,
	}

	if err := s.assetRepo.CreateAsset(ctx, asset); err != nil {
		s.cleanupUploadedFiles(ctx, objectName, thumbnailURL)
		return nil, fmt.Errorf("failed to create asset record: %w", err)
	}

	return asset, nil
}

func (s *AssetService) validateUpload(size int64, contentType string) error {
	if size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}
	if !AllowedImageTypes[contentType] && !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}
	return nil
}

func (s *AssetService) processImage(
	ctx context.Context,
	fileData []byte,
	contentType string,
	isImage bool,
	ext string,
	workspaceID uuid.UUID,
) (width, height *int, thumbnailURL *string, err error) {
	if !isImage {
		return nil, nil, nil, nil
	}

	img, format, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w > MaxImageWidth || h > MaxImageHeight {
		return nil, nil, nil, fmt.Errorf("image dimensions exceed maximum allowed size of %dx%d", MaxImageWidth, MaxImageHeight)
	}

	thumbnailURL, thumbErr := s.createAndUploadThumbnail(ctx, img, format, ext, workspaceID, contentType)
	if thumbErr != nil {
		return nil, nil, nil, thumbErr
	}

	return &w, &h, thumbnailURL, nil
}

func (s *AssetService) createAndUploadThumbnail(
	ctx context.Context,
	img image.Image,
	format, ext string,
	workspaceID uuid.UUID,
	contentType string,
) (*string, error) {
	thumbnail := resize.Thumbnail(ThumbnailWidth, ThumbnailHeight, img, resize.Lanczos3)
	thumbnailName := fmt.Sprintf("%s/%s/thumb_%s%s", workspaceID, time.Now().Format("2006/01"), uuid.New(), ext)

	var thumbnailBuf bytes.Buffer
	var err error
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&thumbnailBuf, thumbnail, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(&thumbnailBuf, thumbnail)
	default:
		err = jpeg.Encode(&thumbnailBuf, thumbnail, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	_, err = s.minioClient.PutObject(
		ctx,
		s.bucketName,
		thumbnailName,
		bytes.NewReader(thumbnailBuf.Bytes()),
		int64(thumbnailBuf.Len()),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	thumbURL := s.getObjectURL(thumbnailName)
	return &thumbURL, nil
}

func (s *AssetService) uploadFile(ctx context.Context, objectName string, fileData []byte, size int64, contentType string) error {
	_, err := s.minioClient.PutObject(ctx, s.bucketName, objectName, bytes.NewReader(fileData), size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (s *AssetService) cleanupUploadedFiles(ctx context.Context, objectName string, thumbnailURL *string) {
	_ = s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if thumbnailURL != nil {
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, *thumbnailURL, minio.RemoveObjectOptions{})
	}
}

// GetAsset retrieves an asset by ID
func (s *AssetService) GetAsset(ctx context.Context, id uuid.UUID) (*models.Asset, error) {
	asset, err := s.assetRepo.GetAssetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return asset, nil
}

// GetWorkspaceAssets retrieves all assets for a workspace
func (s *AssetService) GetWorkspaceAssets(ctx context.Context, workspaceID uuid.UUID) ([]models.Asset, error) {
	assets, err := s.assetRepo.GetAssetsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace assets: %w", err)
	}

	return assets, nil
}

// DeleteAsset soft deletes an asset
func (s *AssetService) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	// Get asset info
	_, err := s.assetRepo.GetAssetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("asset not found: %w", err)
	}

	// Soft delete in database
	if err := s.assetRepo.DeleteAsset(ctx, id); err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	// Note: We don't delete from MinIO immediately to allow for recovery
	// Implement a separate cleanup job for hard deletion

	return nil
}

// CleanupOrphanedAssets finds and deletes assets not referenced by any element
func (s *AssetService) CleanupOrphanedAssets(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	orphanedAssets, err := s.assetRepo.GetOrphanedAssets(ctx, workspaceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get orphaned assets: %w", err)
	}

	count := 0
	for i := range orphanedAssets {
		// Delete from MinIO
		objectName := s.extractObjectName(orphanedAssets[i].URL)
		err := s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			// Log error but continue
			continue
		}

		// Delete thumbnail if exists
		if orphanedAssets[i].ThumbnailURL != nil {
			thumbnailName := s.extractObjectName(*orphanedAssets[i].ThumbnailURL)
			_ = s.minioClient.RemoveObject(ctx, s.bucketName, thumbnailName, minio.RemoveObjectOptions{})
		}

		// Soft delete in database
		if err := s.assetRepo.DeleteAsset(ctx, orphanedAssets[i].ID); err != nil {
			continue
		}

		count++
	}

	return count, nil
}

// Helper functions

func (s *AssetService) getObjectURL(objectName string) string {
	// In production, this should use a CDN URL
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucketName, objectName)
}

func (s *AssetService) extractObjectName(url string) string {
	// Extract object name from full URL
	const urlParts = 2
	parts := strings.SplitN(url, s.bucketName+"/", urlParts)
	if len(parts) == urlParts {
		return parts[1]
	}
	return url
}

// ValidateContentType checks if the content type is allowed
func (s *AssetService) ValidateContentType(contentType string) bool {
	return AllowedImageTypes[contentType]
}
