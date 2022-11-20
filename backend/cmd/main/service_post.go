package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func getFormFile(
	fHdr *multipart.FileHeader,
	MaxFileSizeMb int,
	AllowedFormats []string,
) ([]byte, error) {
	if fHdr.Size > int64(MaxFileSizeMb*1e6) {
		return nil, fmt.Errorf("too big file [%d Mbyte]", fHdr.Size/1e6)
	}

	mimeType := fHdr.Header.Get("Content-Type")
	if !slices.Contains(AllowedFormats, mimeType) {
		return nil, fmt.Errorf("not allowed image format '%s'", mimeType)
	}

	auxiliary, err := fHdr.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		zlog.Err(auxiliary.Close()).Msg("form file closing")
	}()
	return ioutil.ReadAll(auxiliary)
}

func scaleSizes(w, h uint, to uint) (newW, newH uint) {
	if w <= to && h <= to {
		return w, h
	}
	if w == h {
		return to, to
	}
	if w > h {
		return to, h * to / w
	}
	return w * to / h, to
}

func SaveImage(
	imageData []byte,
	filePath string,
	quality uint,
	targetSize uint,
) (err error) {
	imagick.Initialize()
	defer imagick.Terminate()

	wm := imagick.NewMagickWand()
	err = wm.ReadImageBlob(imageData)
	if err != nil {
		zlog.Err(err).Msg("can't read a data of form the file")
		return err
	}
	defer wm.Destroy()

	err = wm.SetImageCompressionQuality(quality)
	if err != nil {
		zlog.Err(err).Msg("can't set a compression quality")
		return err
	}

	err = wm.SetImageFormat("JPEG")
	if err != nil {
		zlog.Err(err).Msg("can't set an image format")
		return err
	}

	err = wm.AutoOrientImage()
	if err != nil {
		zlog.Err(err).Msg("can't set an image orientation")
		return err
	}

	w, h := wm.GetImageWidth(), wm.GetImageHeight()
	if w > targetSize || h > targetSize {
		err = wm.AdaptiveResizeImage(scaleSizes(w, h, targetSize))
		if err != nil {
			zlog.Err(err).Msg("can't resize the image")
			return err
		}
	}
	err = wm.WriteImage(filePath)
	if err != nil {
		zlog.Err(err).Msg("can't write an origin image")
		return err
	}
	return nil
}

func (s *Server) FormProcessing(c *gin.Context) {
	doc := s.dbc.NewDocument()
	err := c.Bind(doc)
	if err != nil {
		zlog.Err(err).Msg("param binding")
		panic(err)
	}

	doc.Dt = time.Now()
	zlog.Debug().Interface("new document ", doc).Msg("")

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		zlog.Err(err).Msg("form file getting")
		panic(err)
	}

	imageData, err := getFormFile(
		fileHeader,
		s.cfg.FormController.MaxFileSize,
		s.cfg.FormController.AllowedFormats,
	)
	if err != nil {
		panic(err)
	}

	if err = SaveImage(
		imageData,
		path.Join(s.cfg.Storage.Path, "o", doc.Id.Hex()+".jpg"),
		s.cfg.FormController.JpegQuality,
		s.cfg.FormController.ShrinkPhotoTo); err != nil {
		panic(err)
	}
	zlog.Debug().Str("image saved", doc.Id.Hex())

	if err = SaveImage(
		imageData,
		path.Join(s.cfg.Storage.Path, "t", doc.Id.Hex()+".jpg"),
		s.cfg.FormController.JpegQuality,
		s.cfg.FormController.ThumbnailSize); err != nil {
		panic(err)
	}
	zlog.Debug().Str("thumbnail saved", doc.Id.Hex())

	if err := s.dbc.Insert(doc); err != nil {
		panic(err)
	}
	zlog.Debug().Str("inserted", doc.Id.Hex())
}
