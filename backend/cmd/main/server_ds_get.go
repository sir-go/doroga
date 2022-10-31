package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetDS(c *gin.Context) {
	id := c.Param("id")
	m, err := regexp.MatchString(`[a-z0-9]{24}`, id)
	if err != nil {
		LOG.Panic(err)
	}
	if !m {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := DBC.RequestsGetById(id)
	if req == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (s *Server) GetDSTar(c *gin.Context) {
	filename := c.Param("filename")
	m, err := regexp.MatchString(`20[0-9]{2}-[01][0-9]-[0-3][0-9]\.tar`, filename)
	if err != nil {
		LOG.Panic(err)
	}
	if !m {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	LOG.Println("get tarball ", filename)

	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		LOG.Panic(err)
	}

	dtString := filename[:10]
	dtm, err := time.ParseInLocation("2006-01-02", dtString, mskLocation)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	LOG.Println("localtime:", dtm)

	nextDay := dtm.AddDate(0, 0, 1)
	reqs := DBC.RequestsFetch(FetchParams{DateBegin: &dtm, DateEnd: &nextDay})

	LOG.Print("found reqs amount:", len(reqs.Records))
	if len(reqs.Records) < 1 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	buf := new(bytes.Buffer)
	tW := tar.NewWriter(buf)
	defer func() {
		err = tW.Close()
		if err != nil {
			LOG.Panic(err)
		}
		c.Header("Content-Type", "application/x-tar")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar\"", dtString))
		_, _ = c.Writer.Write(buf.Bytes())
	}()

	for _, req := range reqs.Records {
		err = tW.WriteHeader(&tar.Header{
			Name:     dtString + "/",
			Mode:     0755,
			ModTime:  req.Dt,
			Typeflag: tar.TypeDir,
		})
		if err != nil {
			LOG.Panic(err)
		}
		err = tW.WriteHeader(&tar.Header{
			Name:     filepath.Join(dtString, req.Id.Hex()+"/"),
			Mode:     0755,
			ModTime:  req.Dt,
			Typeflag: tar.TypeDir,
		})
		if err != nil {
			LOG.Panic(err)
		}

		txtBytes := []byte(req.AsTXT())
		err = tW.WriteHeader(&tar.Header{
			Name:    filepath.Join(dtString, req.Id.Hex(), "data.txt"),
			Size:    int64(len(txtBytes)),
			Mode:    0644,
			ModTime: req.Dt,
		})
		if err != nil {
			LOG.Panic(err)
		}
		_, err = tW.Write(txtBytes)
		if err != nil {
			LOG.Panic(err)
		}

		if imgSrc, err := os.Open(filepath.Join(CFG.Storage.Path, "o", req.Id.Hex()+".jpg")); err == nil {
			stat, err := imgSrc.Stat()
			if err != nil {
				LOG.Panic(err)
			}

			err = tW.WriteHeader(&tar.Header{
				Name:    filepath.Join(dtString, req.Id.Hex(), "photo.jpg"),
				Size:    stat.Size(),
				Mode:    0644,
				ModTime: req.Dt,
			})
			if err != nil {
				LOG.Panic(err)
			}
			_, err = io.Copy(tW, imgSrc)
			if err != nil {
				LOG.Panic(err)
			}

			err = imgSrc.Close()
			if err != nil {
				LOG.Panic(err)
			}
		} else {
			ehSkip(err)
		}
	}
}

func (s *Server) GetDSTarIndex(c *gin.Context) {
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

func (s *Server) GetDSall(c *gin.Context) {
	p := new(FetchParams)
	err := c.BindQuery(p)
	if err != nil {
		LOG.Panic(err)
	}

	reqs := DBC.RequestsFetch(*p)
	c.JSON(http.StatusOK, reqs)
}

func (s *Server) GetDSallPub(c *gin.Context) {
	p := new(FetchParams)
	err := c.BindQuery(p)
	if err != nil {
		LOG.Panic(err)
	}
	p.Fields = "_id,name,bplace,years,vdate,vplace,rang,awards,phdate,info"

	reqs := DBC.RequestsFetch(*p)
	c.JSON(http.StatusOK, reqs)
}
