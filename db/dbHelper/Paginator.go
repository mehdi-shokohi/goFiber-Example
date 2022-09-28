package dbHelper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"goex/config/messages"
)

type Error struct {
	Code    int64         `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type Pagination struct {
	Total       int64       `json:"total,omitempty"`
	PerPage     int64       `json:"per_page,omitempty"`
	CurrentPage int64       `json:"current_page,omitempty"`
	LastPage    *int64      `json:"last_page,omitempty"`
	NextPageUrl *int64      `json:"next_page_url,omitempty"`
	PrevPageUrl *int64      `json:"prev_page_url,omitempty"`
	From        *int64      `json:"from,omitempty"`
	To          *int64      `json:"to,omitempty"`
	Status      bool        `json:"status"`
	Data        interface{} `json:"data"`
	Error       *Error      `json:"error"`
	Sort        *Sort       `json:"-"`
	c           *fiber.Ctx
}

type Sort struct {
	Field string
	Order int8
}

func NewPagination(c *fiber.Ctx) *Pagination {
	m := &Pagination{}
	m.c = c
	s := string(c.Context().QueryArgs().Peek("sort"))
	sorts := strings.Split(s, "|")
	if len(sorts) > 0 && len(sorts[0]) > 0 {
		m.Sort = &Sort{
			Field: sorts[0],
			Order: 1,
		}
		if len(sorts) > 1 {
			if strings.ToLower(sorts[1]) == "desc" {
				m.Sort.Order = -1
			}
		}
	}
	p := string(c.Context().QueryArgs().Peek("page"))
	l := string(c.Context().QueryArgs().Peek("per_page"))

	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		page = 1
	}

	limit, err2 := strconv.ParseInt(l, 10, 64)
	if err2 != nil {
		limit = 0
	}
	m.PerPage = limit
	m.CurrentPage = page
	return m
}

func (p *Pagination) SetLimit(l int64) *Pagination {
	if p.PerPage == 0 || p.PerPage > l {
		p.PerPage = l
	}
	return p
}
func (p *Pagination) SetTotal(t *int64) *Pagination {
	p.Total = *t
	if p.PerPage < 0 {
		p.PerPage = 0
	}
	if p.PerPage == 0 {
		panic(errors.New(messages.InvalidInputForm))
	}
	l := p.Total/p.PerPage + 1
	p.LastPage = &l
	return p
}

func (p *Pagination) DefaultOrder(field string, order int8) *Pagination {
	if p.Sort == nil {
		p.Sort = &Sort{
			Field: field,
			Order: order,
		}
	}
	return p
}

func (p *Pagination) CreateCountOption() *options.CountOptions {
	skip := (p.CurrentPage - 1) * p.PerPage
	return &options.CountOptions{
		Limit:   &p.PerPage,
		MaxTime: nil,
		Skip:    &skip,
	}
}
func (p *Pagination) CreateFindOption() *options.FindOptions {
	skip := (p.CurrentPage - 1) * p.PerPage
	sort := bson.D{{Key: p.Sort.Field, Value: p.Sort.Order}}
	return &options.FindOptions{
		Limit: &p.PerPage,
		Skip:  &skip,
		Sort:  sort,
	}
}
func (p *Pagination) CreateAggregateOptions(pipe *mongo.Pipeline) {
	skip := (p.CurrentPage - 1) * p.PerPage
	sort := bson.D{{Key: p.Sort.Field, Value: p.Sort.Order}}

	*pipe = append(*pipe, bson.D{{"$limit", p.PerPage + skip}})
	*pipe = append(*pipe, bson.D{{"$skip", &skip}})
	*pipe = append(*pipe, bson.D{{"$sort", sort}})

}

//Send New Response Formatter
//Todo: Paginator should not handle response do something about this
func (p *Pagination) Send(data interface{}) {
	p.Status = true
	p.Data = data
	p.c.Response().Header.Set("Content-Type", "application/json")
	p.c.Response().Header.SetStatusCode(200)

	from := (p.CurrentPage-1)*p.PerPage + 1
	to := from + int64(reflect.ValueOf(data).Len()) - 1
	p.From = &from
	p.To = &to
	m := struct {
		Status bool        `json:"status"`
		Data   interface{} `json:"data"`
		Error  *Error      `json:"error"`
	}{
		true,
		p,
		nil}
	v, err := jsoniter.Marshal(m)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprint(p.c, string(v))
	if err != nil {
		fmt.Println("Problem In Response - Internal Error")
	}

}

//Mamali Easy query be thankful
func EasyQuery(c *fiber.Ctx, query bson.D, r Recorder, pagination ...interface{}) error {
	return EasyQueryByDecoder(c, r.GetCollectionName(), query, func(res chan Decoder) ([]interface{}, error) {
		models := make([]interface{}, 0)
		for decoder := range res {
			var m Recorder
			if reflect.ValueOf(r).Kind() == reflect.Ptr {
				m = reflect.New(reflect.TypeOf(r).Elem()).Interface().(Recorder)
			} else {
				m = reflect.New(reflect.TypeOf(r)).Interface().(Recorder)
			}
			err := decoder(m)
			if err != nil {
				return nil, err
			}
			models = append(models, m)
		}
		return models, nil
	}, pagination...)
}

//Mamali Easy query with pain of decoder
func EasyQueryByDecoder(c *fiber.Ctx, collectionName string, query bson.D, f func(chan Decoder) ([]interface{}, error), pagination ...interface{}) error {
	var paginator *Pagination
	var opts *options.FindOptions
	if len(pagination) > 0 {
		paginator = pagination[0].(*Pagination)
		if len(pagination) > 1 {
			opts = pagination[1].(*options.FindOptions)
		} else {
			opts = paginator.CreateFindOption()
		}
	} else {
		paginator = NewPagination(c)
		paginator.
			SetLimit(100).
			DefaultOrder("_id", 1)
		opts = paginator.CreateFindOption()
	}
	var count int64
	res := FindAllGo(c.Context(), collectionName, &query, opts)
	res2 := CountGo(c.Context(), collectionName, &query, &count)

	models, err := f(res)
	if err != nil {
		return err
	}
	err2 := <-res2
	if err2 != nil {
		count = 0
	}
	Count := count
	paginator.SetTotal(&Count)
	paginator.Send(models)
	return nil
}
