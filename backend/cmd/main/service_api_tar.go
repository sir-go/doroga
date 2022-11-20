package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

func parseDate(filename string, location string) (time.Time, error) {
	mskLocation, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	if len(filename) < 10 {
		return time.Time{}, errors.New("filename is too short")
	}
	dtString := filename[:10]
	parsedTime, err := time.ParseInLocation(
		"2006-01-02", dtString, mskLocation)
	return parsedTime, err
}

func getDocumentsOnDate(dbc *DB, dt time.Time) ([]Document, error) {
	nextDay := dt.AddDate(0, 0, 1)
	dbResult, err := dbc.Fetch(FetchParams{DateBegin: &dt, DateEnd: &nextDay})
	if err != nil {
		return nil, err
	}

	zlog.Debug().Int("amount", dbResult.Count).Time("from", dt).Time("to", nextDay).Msg("got documents")
	return dbResult.Docs, nil
}

func tarFile(srcPath string, name string, dt time.Time, tw *tar.Writer) error {
	src, err := os.Open(filepath.Clean(srcPath))
	if err != nil {
		return err
	}

	defer func() {
		if err := src.Close(); err != nil {
			zlog.Err(err).Msg("file closing")
		}
	}()

	stat, err := src.Stat()
	if err != nil {
		return err
	}

	if err = tw.WriteHeader(&tar.Header{
		Name:    name,
		Size:    stat.Size(),
		Mode:    0644,
		ModTime: dt,
	}); err != nil {
		return err
	}

	_, err = io.Copy(tw, src)
	return err
}

func docToTarball(
	doc Document,
	storagePath string,
	rootDir string,
	tw *tar.Writer,
) error {
	txtBytes := []byte(doc.String())
	for _, hdr := range []tar.Header{
		{
			Name:     rootDir + "/",
			Mode:     0755,
			ModTime:  doc.Dt,
			Typeflag: tar.TypeDir,
		},
		{
			Name:     filepath.Join(rootDir, doc.Id.Hex()+"/"),
			Mode:     0755,
			ModTime:  doc.Dt,
			Typeflag: tar.TypeDir,
		},
		{
			Name:    filepath.Join(rootDir, doc.Id.Hex(), "data.txt"),
			Size:    int64(len(txtBytes)),
			Mode:    0644,
			ModTime: doc.Dt,
		},
	} {
		_hdr := hdr
		if err := tw.WriteHeader(&_hdr); err != nil {
			return err
		}
	}

	if _, err := tw.Write(txtBytes); err != nil {
		return err
	}

	return tarFile(
		filepath.Join(storagePath, "o", doc.Id.Hex()+".jpg"),
		filepath.Join(rootDir, doc.Id.Hex(), "photo.jpg"),
		doc.Dt,
		tw,
	)
}

func (s *Server) GetDocumentsTarball(c *gin.Context) {
	filename := c.Param("filename")
	m, err := regexp.MatchString(
		`20[0-9]{2}-[01][0-9]-[0-3][0-9]\.tar`, filename)
	if err != nil {
		panic(err)
	}
	if !m {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	zlog.Debug().Str("filename", filename).Msg("get tar")
	dateString := filename[:10]

	fileDate, err := parseDate(filename, "Europe/Moscow")
	if err != nil {
		panic(err)
	}

	documents, err := getDocumentsOnDate(s.dbc, fileDate)

	buf := new(bytes.Buffer)
	tW := tar.NewWriter(buf)
	defer func() {
		err = tW.Close()
		if err != nil {
			zlog.Err(err).Msg("tar writer closing")
		}
		c.Header("Content-Type", "application/x-tar")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
		c.Header("Content-Disposition", fmt.Sprintf(
			"attachment; filename=\"%s.tar\"", dateString))
		_, _ = c.Writer.Write(buf.Bytes())
	}()

	for _, doc := range documents {
		if err = docToTarball(
			doc, s.cfg.Storage.Path, dateString, tW); err != nil {
			panic(err)
		}
	}
}

func (s *Server) GetTarballsIndex(c *gin.Context) {
	html := `<!DOCTYPE html><html><head>
   	<title>Дорога памяти - архив</title>
   	<meta charset="utf-8">
   	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>p {margin: .5rem;}</style></head>
	<body>
`
	var dayString string
	now := time.Now()
	dateStart := now.AddDate(0, -6, 0)
	for day := now; !day.Before(dateStart); day = day.AddDate(0, 0, -1) {
		dayString = day.Format("2006-01-02")
		html += fmt.Sprintf("<p><a href='%s.tar'>%s</a></p>\n", dayString, dayString)
	}
	html += `
	</body>
	</html>`
	c.String(http.StatusOK, html)
}
