package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

// buildManifestXML создает META-INF/manifest.xml со списком файлов ODS пакета
func buildManifestXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" manifest:version="1.2">
  <manifest:file-entry manifest:full-path="/" manifest:version="1.2" manifest:media-type="application/vnd.oasis.opendocument.spreadsheet"/>
  <manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml"/>
  <manifest:file-entry manifest:full-path="styles.xml" manifest:media-type="text/xml"/>
  <manifest:file-entry manifest:full-path="META-INF/manifest.xml" manifest:media-type="text/xml"/>
</manifest:manifest>`)
}

// buildODSStylesXML создает styles.xml с описанием стилей документа и параметров страницы
func buildODSStylesXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<office:document-styles xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0"
  xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0"
  xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0"
  xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0"
  xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0"
  xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0"
  office:version="1.2">
  <office:font-face-decls>
    <style:font-face style:name="Calibri" svg:font-family="Calibri, 'Liberation Sans', sans-serif"/>
  </office:font-face-decls>
  <office:styles>
    <style:default-style style:family="table-cell">
      <style:text-properties style:font-name="Calibri" fo:font-size="11pt" fo:language="ru" fo:country="RU"/>
    </style:default-style>
  </office:styles>
  <office:automatic-styles>
    <style:page-layout style:name="pm1">
      <style:page-layout-properties fo:page-width="21.0cm" fo:page-height="29.7cm" style:print-orientation="portrait"
        fo:margin-top="1.0cm" fo:margin-bottom="1.0cm" fo:margin-left="1.0cm" fo:margin-right="1.0cm"
        style:scale-to="100%"/>
    </style:page-layout>
  </office:automatic-styles>
  <office:master-styles>
    <style:master-page style:name="Default" style:page-layout-name="pm1"/>
  </office:master-styles>
</office:document-styles>`)
}

// writeZipEntry записывает файл в ZIP-архив с указанным методом сжатия
func writeZipEntry(zw *zip.Writer, name string, method uint16, data []byte) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: method,
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("ошибка создания записи zip %s: %w", name, err)
	}
	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("ошибка записи данных zip %s: %w", name, err)
	}
	return nil
}

