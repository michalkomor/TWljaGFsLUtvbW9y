package item

type Getter interface {
	GetAll() []Item
}

type Adder interface {
	Add(item Item)
}

type Item struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Interval int       `json:"interval"`
	History  []History `json:"-"`
}

type History struct {
	Response  string  `json:"response"`
	Duration  float64 `json:"duration"`
	CreatedAt int64   `json:"created_at"`
}

type Collection struct {
	Items []Item
}

func New() *Collection {
	return &Collection{}
}

func (c *Collection) Add(item Item) {
	c.Items = append(c.Items, item)
}

func (c *Collection) GetAll() []Item {
	return c.Items
}

func (c *Collection) Delete(i int) {
	c.Items = append(c.Items[:i], c.Items[i+1:]...)
}
