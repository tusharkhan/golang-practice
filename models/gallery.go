package models

import (
	"course/helper"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Gallery struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type Image struct {
	Path      string
	GalleryId int
	FileName  string
}

type GalleryService struct {
	DB *sql.DB

	ImageDirectory string
}

func (gs *GalleryService) Create(title string, user_id int) (*Gallery, error) {
	var gall Gallery = Gallery{
		Title:     title,
		UserID:    user_id,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	var galleryInsertQuery string = "INSERT INTO galleries (user_id, title) VALUES ($1, $2) RETURNING id"
	galleryInsertQueryError := gs.DB.QueryRow(galleryInsertQuery, user_id, title).Scan(&gall.ID)

	if galleryInsertQueryError != nil {
		return nil, fmt.Errorf("error in creating gallery %w", galleryInsertQueryError)
	}

	return &gall, nil
}

func (gs *GalleryService) Show(id int) (*Gallery, error) {
	var gall Gallery = Gallery{
		ID: id,
	}

	var singleGalleryQuery string = "SELECT user_id,  title, created_at FROM galleries WHERE id=$1"
	singleGalleryQueryError := gs.DB.QueryRow(singleGalleryQuery, gall.ID).Scan(&gall.UserID, &gall.Title, &gall.CreatedAt)

	if singleGalleryQueryError != nil {
		if errors.Is(singleGalleryQueryError, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("error in getting single gallery %w", singleGalleryQueryError)
	}

	return &gall, nil
}

func (gs *GalleryService) List() ([]Gallery, error) {
	var listQuery string = `SELECT galleries.id, galleries.title, galleries.user_id, galleries.created_at, users.name FROM galleries 
	LEFT JOIN users ON users.id = galleries.user_id`

	rows, listQueryError := gs.DB.Query(listQuery)

	if listQueryError != nil {
		return nil, fmt.Errorf("error in getting gallery list %w", listQueryError)
	}

	var galleries []Gallery

	for rows.Next() {
		var gal Gallery
		err := rows.Scan(&gal.ID, &gal.Title, &gal.UserID, &gal.CreatedAt, &gal.User.Name)

		if err != nil {
			return nil, fmt.Errorf("error in fetching gallery %w", err)
		}

		galleries = append(galleries, gal)
	}

	return galleries, nil
}

func (gs *GalleryService) GetByUser(user_id int) ([]Gallery, error) {
	var getQueryString string = "SELECT  id, title, created_at FROM galleries WHERE user_id = $1"

	rows, queryError := gs.DB.Query(getQueryString, user_id)

	if queryError != nil {
		if errors.Is(queryError, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("error in getting gallery %w", queryError)
	}

	var galleries []Gallery

	for rows.Next() {
		gall := Gallery{
			UserID: user_id,
		}

		queryError = rows.Scan(&gall.ID, &gall.Title, &gall.CreatedAt)

		if queryError != nil {
			return nil, fmt.Errorf("error in scaning data %w", queryError)
		}

		galleries = append(galleries, gall)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in getting row %w", rows.Err())
	}

	return galleries, nil
}

func (gs *GalleryService) Update(gall *Gallery) error {
	var updateQuery string = "UPDATE galleries SET title = $1 WHERE id = $2"
	_, updateQueryError := gs.DB.Exec(updateQuery, gall.Title, gall.ID)

	if updateQueryError != nil {
		return fmt.Errorf("error in updating galleries %w", updateQueryError)
	}

	return nil
}

func (gs *GalleryService) Delete(id int) error {
	var galleryDeleteQuery string = "DELETE FROM galleries WHERE id=$1"

	_, deleteError := gs.DB.Exec(galleryDeleteQuery, id)

	if deleteError != nil {
		return fmt.Errorf("error in deleting gallery %w", deleteError)
	}

	return nil
}

func (gs *GalleryService) GalleryDire(id int) string {
	var imgDir string = gs.ImageDirectory
	if imgDir == "" {
		imgDir = "images"
	}

	return filepath.Join(imgDir, fmt.Sprintf("gallery-%d", id))
}

func (gs *GalleryService) Images(galleryId int) ([]Image, error) {
	globPattern := filepath.Join(gs.GalleryDire(galleryId), "*")
	allFiles, fileFetchError := filepath.Glob(globPattern)

	if fileFetchError != nil {
		return nil, fmt.Errorf("error in fetching file data %w", fileFetchError)
	}

	var imagePaths []Image

	for _, fil := range allFiles {
		if helper.HasExtension(fil, []string{".png", ".jpg", ".jpeg", "gif"}) {
			imagePaths = append(imagePaths, Image{
				Path:      fil,
				GalleryId: galleryId,
				FileName:  filepath.Base(fil),
			})
		}
	}

	return imagePaths, nil
}

func (gs *GalleryService) UploadImage(galleryId int, fileNmae string, content io.Reader) (string, error) {
	var galleryDir string = gs.GalleryDire(galleryId)
	directoryCreatingError := os.MkdirAll(galleryDir, 0775)

	if directoryCreatingError != nil {
		return "", fmt.Errorf("error in creating gallery directory %w", directoryCreatingError)
	}

	var imagePath string = filepath.Join(galleryDir, fileNmae)

	createImage, createImageError := os.Create(imagePath)

	if createImageError != nil {
		return "", fmt.Errorf("error in creating image %w", createImageError)
	}

	defer createImage.Close()

	_, imageCopyError := io.Copy(createImage, content)

	if imageCopyError != nil {
		return "", fmt.Errorf("error in copying imaeg to directory %w", imageCopyError)
	}

	return imagePath, nil
}
