package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func getFormFile(fHdr *multipart.FileHeader) ([]byte, error) {
	if fHdr.Size > CFG.ReqWebForm.MaxFileSize*1e6 {
		return nil, fmt.Errorf("too bug file [%d Mbyte]", fHdr.Size/1e6)
	}

	mimeType := fHdr.Header.Get("Content-Type")
	if !StringSliceContains(CFG.ReqWebForm.AllowedFormats, mimeType) {
		return nil, fmt.Errorf("not allowed image format '%s'", mimeType)
	}

	auxiliar, err := fHdr.Open()
	if err != nil {
		return nil, err
	}
	defer func() { ehSkip(auxiliar.Close()) }()
	return ioutil.ReadAll(auxiliar)
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

func SaveImage(fh *multipart.FileHeader, reqId string) {
	imagick.Initialize()
	defer imagick.Terminate()
	imageBytes, err := getFormFile(fh)
	if err != nil {
		LOG.Panic("can't read the form file", err)
	}

	wm := imagick.NewMagickWand()
	err = wm.ReadImageBlob(imageBytes)
	if err != nil {
		LOG.Panic("can't read a data of form the file", err)
	}
	defer wm.Destroy()

	err = wm.SetImageCompressionQuality(CFG.ReqWebForm.JpegQuality)
	if err != nil {
		LOG.Panic("can't set a compression quality", err)
	}

	err = wm.SetImageFormat("JPEG")
	if err != nil {
		LOG.Panic("can't set an image format", err)
	}

	err = wm.AutoOrientImage()
	if err != nil {
		LOG.Panic("can't set an image orientation", err)
	}

	w, h := wm.GetImageWidth(), wm.GetImageHeight()
	if w > CFG.ReqWebForm.ShrinkPhotoTo || h > CFG.ReqWebForm.ShrinkPhotoTo {
		err = wm.AdaptiveResizeImage(scaleSizes(w, h, CFG.ReqWebForm.ShrinkPhotoTo))
		if err != nil {
			LOG.Panic("can't resize the image", err)
		}
	}
	err = wm.WriteImage(path.Join(CFG.Storage.Path, "o", reqId+".jpg"))
	if err != nil {
		LOG.Panic("can't write an origin image", err)
	}

	w, h = wm.GetImageWidth(), wm.GetImageHeight()
	if w > CFG.ReqWebForm.ThumbnailSize || h > CFG.ReqWebForm.ThumbnailSize {
		err = wm.ThumbnailImage(scaleSizes(w, h, CFG.ReqWebForm.ThumbnailSize))
		if err != nil {
			LOG.Panic("can't make a thumb image", err)
		}
	}
	err = wm.WriteImage(path.Join(CFG.Storage.Path, "t", reqId+".jpg"))
	if err != nil {
		LOG.Panic("can't write a thumbnail image", err)
	}
}

func (s *Server) FormProcessing(c *gin.Context) {
	dbReq := DBC.NewReq()
	err := c.Bind(dbReq)
	if err != nil {
		LOG.Panic(err)
	}

	dbReq.Dt = time.Now()
	LOG.Println(dbReq)

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		LOG.Panic(err)
	}

	SaveImage(fileHeader, dbReq.Id.Hex())
	LOG.Println(dbReq.Id.Hex(), ": images saved")

	if err := DBC.RequestInsert(dbReq); err != nil {
		c.Status(http.StatusInternalServerError)
		ehSkip(err)
		return
	}
	LOG.Println(dbReq.Id.Hex(), ": inserted")
}
