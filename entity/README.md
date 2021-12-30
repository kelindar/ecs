# Entity

This package contains various **entities** of the game, things such as players, items, monsters etc. An entity represents something that has a set of **components** (i.e. columns) and is expressed as a **view** over a row that is lazily evaluated by various **systems**.

For example, a mobile entity takes a pointer over a specific row and allows us to read/write columns by simply accessing the properties instead of dealing with the internals of the data storage layer. In the example below, we have `Location() tile.Point` and `SetLocation(tile.Point)` methods on our mobile entity, which contain the necessary code to read/write but the entity itself does not have any state besides a `Cursor`. This allows us to lazily evaluate and only read the necessary data when accessing it.

```go
// Mobile represents a view on a current mobile row
type Mobile struct {
	row *column.Cursor
}

// Location reads the current location
func (m *Mobile) Location() tile.Point {
    // read the "location" column
}

// SetLocation writes the current location
func (m *Mobile) SetLocation(v tile.Point) {
	// write the "location" column
}
```
