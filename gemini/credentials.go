package gemini

type Credentials struct {
	APIKey string
}

func (c *Credentials) ApplyToModel(model *Model) {
	model.setCredentials(c)
}
