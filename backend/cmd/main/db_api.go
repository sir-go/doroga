package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	zlog "github.com/rs/zerolog/log"
)

type Document struct {
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

func (r *Document) String() (res string) {
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
	Docs  []Document `json:"records"`
	Count int        `json:"count"`
	Total int        `json:"total"`
}

func (db *DB) initSession() (*mgo.Collection, func()) {
	session := db.session.Copy()
	return session.DB("").C(db.collectionName), session.Close
}

func (db *DB) GetAll() (*FetchResult, error) {
	c, closer := db.initSession()
	defer closer()

	var items []Document

	q := c.Find(nil)

	count, err := q.Count()
	if err != nil {
		return nil, err
	}

	if err = q.All(&items); err != nil {
		return nil, err
	}

	return &FetchResult{Docs: items, Count: count}, nil
}

func (db *DB) Fetch(p FetchParams) (*FetchResult, error) {
	c, closer := db.initSession()
	defer closer()

	items := make([]Document, 0)

	var query, qSearch, qDate bson.M

	if p.Search != "" {
		searchCriteria := make([]bson.M, 0)
		if p.SearchFields == "" || p.SearchFields == "name" {
			searchCriteria = append(searchCriteria, bson.M{
				"name": bson.RegEx{Pattern: p.Search, Options: "i"}})
		} else {
			for _, fieldName := range strings.Split(p.SearchFields, ",") {
				searchCriteria = append(searchCriteria, bson.M{
					fieldName: bson.RegEx{Pattern: p.Search, Options: "i"}})
			}
		}
		qSearch = bson.M{"$or": searchCriteria}
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
		zlog.Err(err).Msg("db get query records count")
		return nil, err
	}

	total, err := c.Count()
	if err != nil {
		zlog.Err(err).Msg("db get total records count")
		return nil, err
	}

	if p.SortField != "" {
		q = q.Sort(p.SortField)
	}

	if p.Limit > 0 {
		q = q.Limit(p.Limit)
	}

	if p.Offset > 0 {
		zlog.Debug().Int("offset", p.Offset).Msg("")
		if p.Offset > count {
			zlog.Debug().Msg("> total count then offset will by ignored")
		} else {
			q = q.Skip(p.Offset)
		}
	}

	if err = q.All(&items); err != nil {
		zlog.Err(err).Msg("db get all items")
		return nil, err
	}

	return &FetchResult{Docs: items, Count: count, Total: total}, nil
}

func (db *DB) GetTotal() (int, error) {
	c, closer := db.initSession()
	defer closer()

	q, err := c.Count()
	if err != nil {
		zlog.Err(err).Msg("get records total amount")
		return 0, err
	}
	return q, nil
}

func (db *DB) ExistsById(id string) (bool, error) {
	c, closer := db.initSession()
	defer closer()

	if !bson.IsObjectIdHex(id) {
		err := fmt.Errorf("id is not a valid mongodb object id")
		zlog.Err(err).Msg("exists by id")
		return false, err
	}

	count, err := c.FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (db *DB) GetAllIds() ([]string, error) {
	c, closer := db.initSession()
	defer closer()

	var (
		Records []Document
		res     []string
	)

	err := c.Find(bson.M{}).Select(bson.M{"_id": 1}).All(&Records)
	if err != nil {
		return nil, err
	}
	for _, rec := range Records {
		res = append(res, rec.Id.Hex())
	}

	return res, nil
}

func (db *DB) GetById(id string) (*Document, error) {
	c, closer := db.initSession()
	defer closer()

	var item Document
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&item); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	if item.Id.Valid() {
		return &item, nil
	}
	return nil, fmt.Errorf("record id is not valid")
}

func (db *DB) GetByUUID(uuid string) (*Document, error) {
	c, closer := db.initSession()
	defer closer()

	var item Document
	if err := c.Find(bson.M{"uuid": uuid}).One(&item); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	if item.Id.Valid() {
		return &item, nil
	}
	return nil, fmt.Errorf("record id is not valid")
}

func (db *DB) NewDocument() *Document {
	return &Document{Id: bson.NewObjectId()}
}

func (db *DB) Insert(item *Document) error {
	c, closer := db.initSession()
	defer closer()

	if !item.Id.Valid() {
		item.Id = bson.NewObjectId()
	}

	return c.Insert(item)
}
