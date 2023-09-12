package paginator

import (
	"errors"
	"math"

	"gorm.io/gorm"
)

type FuncType int

const maxInt = 9223372036854775807

const (
	_       FuncType = iota
	Find             //gorm Find
	Related          //gorm Related
)

type Pagination struct {
	Total       int64       `json:"total"`
	PerPage     int         `json:"per_page"`
	PrePage     int         `json:"pre_page"`
	CurrentPage int         `json:"current_page"`
	NextPage    int         `json:"next_page"`
	LastPage    int         `json:"last_page"`
	Rows        interface{} `json:"rows"`
}

type Params struct {
	GQuery     *gorm.DB //gorm's query,like gorm.db.where("xxx")...etc
	PerPage    int
	Page       int
	FuncName   FuncType //gorm's query function:find or related
	RelatedKey string   //related forigen key
}

// Paginate takes params
func Paginate(param Params, result interface{}) (*Pagination, error) {
	var pagination Pagination

	gquery := param.GQuery

	//calc total
	switch param.FuncName {
	case 0:
		if err := gquery.Model(result).Count(&pagination.Total).Error; err != nil {
			return nil, err
		}
	case Find:
		if err := gquery.Model(result).Count(&pagination.Total).Error; err != nil {
			return nil, err
		}
	case Related:
		pagination.Total = gquery.Association(param.RelatedKey).Count()
	default:
		return nil, errors.New("func type not supported")

	}

	pagination.PerPage = param.PerPage
	if pagination.PerPage <= 0 {
		pagination.PerPage = 10
	}
	pagination.LastPage = int(math.Ceil(float64(pagination.Total) / float64(pagination.PerPage)))
	if param.Page < 1 {
		param.Page = 1
	} else if param.Page > pagination.LastPage {
		param.Page = pagination.LastPage
	}
	pagination.CurrentPage = param.Page
	if (param.Page - 1) <= 0 {
		pagination.PrePage = 1
	} else {
		pagination.PrePage = param.Page - 1
	}
	if param.Page < maxInt {
		if (param.Page + 1) >= pagination.LastPage {
			pagination.NextPage = pagination.LastPage
		} else {
			pagination.NextPage = param.Page + 1
		}
	} else {
		pagination.NextPage = maxInt
	}

	offset := (pagination.CurrentPage - 1) * pagination.PerPage
	//query data
	switch param.FuncName {
	case 0:
		if err := gquery.Limit(pagination.PerPage).Offset(offset).Find(result).Error; err != nil {
			return nil, err
		}
	case Find:
		if err := gquery.Limit(pagination.PerPage).Offset(offset).Find(result).Error; err != nil {
			return nil, err
		}
	case Related:
		if err := gquery.Limit(pagination.PerPage).Offset(offset).Association(param.RelatedKey).Find(result); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("func type not supported")

	}

	pagination.Rows = result

	return &pagination, nil
}
