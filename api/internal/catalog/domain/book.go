package domain

type Testament string

const (
	TestamentOld Testament = "old"
	TestamentNew Testament = "new"
)

type Book struct {
	id           string
	order        int
	name         string
	testament    Testament
	chapterCount int
}

type NewBookParams struct {
	ID           string
	Order        int
	Name         string
	Testament    Testament
	ChapterCount int
}

func NewBook(params NewBookParams) (Book, error) {
	if params.ID == "" {
		return Book{}, ErrInvalidBookID
	}
	if params.Order <= 0 {
		return Book{}, ErrInvalidBookOrder
	}
	if params.Name == "" {
		return Book{}, ErrInvalidBookName
	}
	if params.Testament != TestamentOld && params.Testament != TestamentNew {
		return Book{}, ErrInvalidTestament
	}
	if params.ChapterCount <= 0 {
		return Book{}, ErrInvalidChapterCount
	}
	return Book{id: params.ID, order: params.Order, name: params.Name, testament: params.Testament, chapterCount: params.ChapterCount}, nil
}

func (b Book) ID() string           { return b.id }
func (b Book) Order() int           { return b.order }
func (b Book) Name() string         { return b.name }
func (b Book) Testament() Testament { return b.testament }
func (b Book) ChapterCount() int    { return b.chapterCount }
