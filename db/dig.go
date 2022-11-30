package db

// Dig exported
type Dig struct {
	Table   Table
	Primary Params
}

func NewDig(t Table, p Params) Dig {

	return Dig{Table: t, Primary: p}
}

func dig(qo *QueryOptions) error {

	for _, v := range qo.Table.Dig(qo.Dig...) {
		qp := make(QParams)
		qp["dig"] = qo.Dig
		dqo := QueryOptionsFactory(v.Table, qo.UID, qp, UParams(v.Primary))
		if _, err := Find(dqo); err != nil {
			return err
		}
	}
	return nil
}
