# Systems

This package contains various game **systems** that are executed periodically and process a set of **components** (i.e. columns) for a set of **entities** (i.e. players, items, monsters). Systems access data using columnar **queries** which allow us to filter only the rows that the system can process.

For example, consider a _movement system_ that adjust both `location` and `movement action` components of a `mobile object` entity. Such system needs to first filter out entities that haven't moved, which is done using a bitmap index that checks whether the `movement action` is not empty, then the system processes all of the matching entities by updating the movement accordingly. If a lot of entities move, this is very efficient since cache misses are reduced and the movement logic is neatly contained within the system that performs the given behavior.

```go
// Update is called periodically on the movement system
func (s *System) Update(dt time.Duration) error {
	return s.mobs.Query(func(txn *column.Txn) error {
		return txn.With("moving").Range("move", func(row column.Cursor) {
			mobile := entity.MobileAt(&row)
			movement := m.Movement()
            location := m.Location()

            // Update the movement vector and store it
            movement, moved = movement.Update(dt)
            m.SetMovement(movement)
            if !moved {
                return false // not moved
            }

            // Try to move and check whether the location is within map bounds
            location = location.Move(movement.Direction())
            if !location.WithinSize(s.grid.Size){
                return false // out of map bounds
            }

            // Update the current location
            m.SetLocation(location)
		})
	})
}
```
