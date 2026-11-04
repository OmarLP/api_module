package meta

import (
	"os"
	"strconv"
)

type Meta struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

func New(page, perPage, total int64) (*Meta, error) {
	// si las filas que va a mostrar es menor o igual a 0 le asignamos por defeco
	if perPage <= 0 {
		var err error
		perPage, err = strconv.ParseInt(os.Getenv("PAGINATOR_LIMIT_DEFAULT"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// cantidad de paginas en la respuesta
	pageCount := int64(0)
	if total >= 0 {
		pageCount = int64(total+perPage-1) / perPage
		// menos pag de las que pasó por parámetro le asigo la últma
		if page > pageCount {
			page = pageCount
		}
	}

	// si la pag es 0 o menor, asignar 1
	if page < 1 {
		page = 1
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		PageCount:  pageCount,
	}, nil
}

// desde que número de fila
func (p *Meta) Offset() int64 {
	return (p.Page - 1) * p.PerPage
}

// hasta que número de fila
func (p *Meta) Limit() int64 {
	return p.PerPage
}
