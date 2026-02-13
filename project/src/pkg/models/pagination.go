package models

type PaginationResult[T any] struct {
	Page         int
	ItemsPerPage int
	TotalItems   int
	TotalPages   int
	Items        []T
}

type PaginationInput[S ~string] struct {
	Page  int
	Limit int
	Sort  *S
}

type PaginationMeta struct {
	Offset int
	Limit  int
}

func NormalizePagination(page, limit int) PaginationMeta {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return PaginationMeta{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func ToAnySlice[T any](items []T) []any {
	out := make([]any, len(items))
	for i := range items {
		out[i] = items[i]
	}
	return out
}

func MapPagination[T any, D any](
	pr PaginationResult[T],
	mapFn func(T) D,
) PaginationResult[D] {
	items := make([]D, len(pr.Items))
	for i, it := range pr.Items {
		items[i] = mapFn(it)
	}

	return PaginationResult[D]{
		Page:         pr.Page,
		ItemsPerPage: pr.ItemsPerPage,
		TotalItems:   pr.TotalItems,
		TotalPages:   pr.TotalPages,
		Items:        items,
	}
}

func MapDomainSliceToSliceDTO[T any, D any](
	domain []T,
	mapFn func(T) D,
) []D {
	items := make([]D, len(domain))
	for i, it := range domain {
		items[i] = mapFn(it)
	}
	return items
}
