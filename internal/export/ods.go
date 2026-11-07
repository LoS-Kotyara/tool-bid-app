package export

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"tool-bid-app/internal/model"
)

// ExportODS генерирует файл OpenDocument Spreadsheet (.ods) автономно
// через стандартные библиотеки archive/zip и encoding/xml без сторонних зависимостей.
func ExportODS(doc *model.Document, outputPath string) error {
	cleanPath := filepath.Clean(outputPath)
	if err := os.MkdirAll(filepath.Dir(cleanPath), 0755); err != nil {
		return fmt.Errorf("ошибка создания каталога назначения: %w", err)
	}

	f, err := os.Create(cleanPath)
	if err != nil {
		return fmt.Errorf("ошибка создания ODS файла: %w", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	// 1. mimetype: ОБЯЗАН быть первым файлом, сохраненным БЕЗ компрессии (Store)
	mimetypeData := []byte("application/vnd.oasis.opendocument.spreadsheet")
	if err := writeZipEntry(zw, "mimetype", zip.Store, mimetypeData); err != nil {
		return err
	}

	// 2. META-INF/manifest.xml
	manifestData := buildManifestXML()
	if err := writeZipEntry(zw, "META-INF/manifest.xml", zip.Deflate, manifestData); err != nil {
		return err
	}

	// 3. styles.xml
	stylesData := buildODSStylesXML()
	if err := writeZipEntry(zw, "styles.xml", zip.Deflate, stylesData); err != nil {
		return err
	}

	// 4. content.xml
	contentData := buildODSContentXML(doc)
	if err := writeZipEntry(zw, "content.xml", zip.Deflate, contentData); err != nil {
		return err
	}

	return nil
}
