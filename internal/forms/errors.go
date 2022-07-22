package forms

type errors map[string][]string

//Adds an error message for given form field
func (e errors) Add(field, messages string) {
	e[field] = append(e[field], messages)
}

//Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}