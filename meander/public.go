package meander

// Facade is
type Facade interface {
	Public() interface{}
}

// Public is
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
