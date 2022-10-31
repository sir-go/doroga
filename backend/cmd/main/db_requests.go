package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DBreq struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Uuid       string        `json:"uuid,omitempty" bson:",omitempty"`
	Dt         time.Time     `json:"dt,omitempty" bson:",omitempty"`
	Name       string        `json:"name,omitempty" form:"name" binding:"required,lte=4000" bson:",omitempty" txt:"имя"`
	Bplace     string        `json:"bplace,omitempty" form:"bplace" binding:"lte=4000" bson:",omitempty" txt:"место рождения"`
	Years      string        `json:"years,omitempty" form:"years" binding:"lte=4000" bson:",omitempty" txt:"годы жизни"`
	Vdate      string        `json:"vdate,omitempty" form:"vdate" binding:"lte=4000" bson:",omitempty" txt:"дата призыва"`
	Vplace     string        `json:"vplace,omitempty" form:"vplace" binding:"lte=4000" bson:",omitempty" txt:"пункт призыва"`
	Rang       string        `json:"rang,omitempty" form:"rang" binding:"lte=4000" bson:",omitempty" txt:"воинское звание"`
	Awards     string        `json:"awards,omitempty" form:"awards" binding:"lte=4000" bson:",omitempty" txt:"награды"`
	PhDate     string        `json:"ph_date,omitempty" form:"phDate" binding:"lte=4000" bson:",omitempty" txt:"когда сделано фото"`
	Info       string        `json:"info,omitempty" form:"info" binding:"lte=100000" bson:",omitempty" txt:"доп. сведения"`
	SenderName string        `json:"sender_name,omitempty" form:"sender_name" binding:"required,lte=4000" bson:",omitempty" txt:"контакт фио"`
	Phone      string        `json:"phone,omitempty" form:"phone" binding:"required,lte=4000" bson:",omitempty" txt:"контакт телефон"`
	VKpost     string        `json:"vkpost,omitempty" bson:",omitempty"`
}

func (r *DBreq) AsTXT() (res string) {
	s := reflect.ValueOf(r).Elem()
	t := reflect.TypeOf(*r)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("txt")
		fVal := s.Field(i)
		if tag == "" || fVal.Interface() == "" {
			continue
		}
		res += fmt.Sprintf("[%s]\r\n%s\r\n\r\n", tag, fVal.Interface())
	}
	return
}

type FetchParams struct {
	DateBegin       *time.Time `form:"date_begin"`
	DateEnd         *time.Time `form:"date_end"`
	Search          string     `form:"search"`
	SearchFields    string     `form:"search_fields"`
	UnpublishedOnly bool       `form:"unpublished_only"`
	SortField       string     `form:"sort"`
	Limit           int        `form:"limit"`
	Offset          int        `form:"offset"`
	Fields          string     `form:"fields"`
}

type FetchResult struct {
	Records []DBreq `json:"records"`
	Count   int     `json:"count"`
	Total   int     `json:"total"`
}

func (db *DB) RequestsGetAll() *FetchResult {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	var items []DBreq

	q := c.Find(nil)

	count, err := q.Count()
	if err != nil {
		LOG.Panic(err)
	}

	if err = q.All(&items); err != nil {
		LOG.Panic(err)
	}
	return &FetchResult{Records: items, Count: count}
}

func (db *DB) RequestsFetch(p FetchParams) *FetchResult {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	items := make([]DBreq, 0)

	var query, qSearch, qDate bson.M

	if p.Search != "" {
		searchCriterias := make([]bson.M, 0)
		if p.SearchFields == "" || p.SearchFields == "name" {
			searchCriterias = append(searchCriterias, bson.M{
				"name": bson.RegEx{Pattern: p.Search, Options: "i"}})
		} else {
			for _, fname := range strings.Split(p.SearchFields, ",") {
				searchCriterias = append(searchCriterias, bson.M{
					fname: bson.RegEx{Pattern: p.Search, Options: "i"}})
			}
		}
		qSearch = bson.M{"$or": searchCriterias}
	}

	if p.DateBegin != nil || p.DateEnd != nil {
		qDate = bson.M{"dt": bson.M{}}
	}

	if p.DateBegin != nil {
		qDate["dt"].(bson.M)["$gte"] = p.DateBegin
	}
	if p.DateEnd != nil {
		qDate["dt"].(bson.M)["$lt"] = p.DateEnd
	}

	query = bson.M{"$and": []bson.M{qSearch, qDate}}
	q := c.Find(query)

	if len(p.Fields) > 0 {
		fields := bson.M{}
		for _, f := range strings.Split(p.Fields, ",") {
			fields[f] = 1
		}
		q = q.Select(fields)
	}

	q = q.Collation(&mgo.Collation{Locale: "ru"})

	count, err := q.Count()
	if err != nil {
		LOG.Panic(err)
	}

	total, err := c.Count()
	if err != nil {
		LOG.Panic(err)
	}

	if p.SortField != "" {
		q = q.Sort(p.SortField)
	}

	if p.Limit > 0 {
		q = q.Limit(p.Limit)
	}

	if p.Offset > 0 {
		if p.Offset > count {
			LOG.Println("query Offset: ", p.Offset, ", it > total count then offset will by ignored")
		} else {
			LOG.Println("query Offset: ", p.Offset)
			q = q.Skip(p.Offset)
		}
	}

	err = q.All(&items)
	if err != nil {
		LOG.Panic(err)
	}
	return &FetchResult{Records: items, Count: count, Total: total}
}

func (db *DB) RequestsGetTotal() int {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	q, err := c.Count()
	if err != nil {
		LOG.Panic(err)
	}
	return q
}

func (db *DB) RequestsExistsById(id string) bool {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	if !bson.IsObjectIdHex(id) {
		return false
	}

	count, err := c.FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		if ehIsNotFound(err) {
			return false
		}
		if err != nil {
			LOG.Panic(err)
		}
	}
	if count > 0 {
		return true
	}
	return false
}

func (db *DB) RequestsGetAllIds() []string {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	var (
		Records []DBreq
		res     []string
	)

	err := c.Find(bson.M{}).Select(bson.M{"_id": 1}).All(&Records)
	if err != nil {
		LOG.Panic(err)
	}
	for _, rec := range Records {
		res = append(res, rec.Id.Hex())
	}

	return res
}

func (db *DB) RequestsGetById(id string) *DBreq {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	var item DBreq
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&item); err != nil {
		if ehIsNotFound(err) {
			return nil
		}
		if err != nil {
			LOG.Panic(err)
		}
	}

	if item.Id.Valid() {
		return &item
	}
	return nil
}

func (db *DB) RequestsGetByUUID(uuid string) *DBreq {
	session := db.session.Copy()
	defer session.Close()
	c := session.DB("").C(CFG.Db.Collection)

	var item DBreq
	if err := c.Find(bson.M{"uuid": uuid}).One(&item); err != nil {
		if ehIsNotFound(err) {
			return nil
		}
		if err != nil {
			LOG.Panic(err)
		}
	}

	if item.Id.Valid() {
		return &item
	}
	return nil
}

func (db *DB) NewReq() *DBreq {
	return &DBreq{Id: bson.NewObjectId()}
}

func (db *DB) RequestInsert(item *DBreq) error {
	session := db.session.Copy()
	defer session.Close()
	if !item.Id.Valid() {
		item.Id = bson.NewObjectId()
	}
	return session.DB("").C(CFG.Db.Collection).Insert(item)
}

func (db *DB) RequestRemove(id string) error {
	session := db.session.Copy()
	defer session.Close()
	if err := session.DB("").C(CFG.Db.Collection).RemoveId(bson.ObjectIdHex(id)); err != nil {
		if ehIsNotFound(err) {
			LOG.Println(id, err)
			return nil
		}
		return err
	}
	ehSkip(os.Remove(path.Join(CFG.Storage.Path, "o", id+".jpg")))
	ehSkip(os.Remove(path.Join(CFG.Storage.Path, "t", id+".jpg")))
	return nil
}
