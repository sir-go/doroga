package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetDocuments(c *gin.Context) {
	id := c.Param("id")
	m, err := regexp.MatchString(`[a-z0-9]{24}`, id)
	if err != nil {
		panic(err)
	}
	if !m {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req, err := s.dbc.GetById(id)
	if err != nil {
		panic(err)
	}
	if req == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (s *Server) GetDocumentsAll(c *gin.Context) {
	p := new(FetchParams)
	if err := c.BindQuery(p); err != nil {
		panic(err)
	}

	reqs, err := s.dbc.Fetch(*p)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, reqs)
}

func (s *Server) GetAllDocumentsPub(c *gin.Context) {
	p := new(FetchParams)
	if err := c.BindQuery(p); err != nil {
		panic(err)
	}
	p.Fields = "_id,name,bplace,years,vdate,vplace,rang,awards,phdate,info"

	reqs, err := s.dbc.Fetch(*p)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, reqs)
}
