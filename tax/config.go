package tax

type Config struct {
	Name  string  `db:"name"`
	Value float64 `db:"value"`
}
