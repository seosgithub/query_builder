package query_builder

type QueryMixin struct {
	*QueryBuilder
}

type QueryBuilder struct {
	queries []interface{}

	QueryRunner
}

// Called inside mixins to add a query (representation) to the stack
func (qb *QueryBuilder) Push(query interface{}) {
	qb.queries = append(qb.queries, query)
}

type QueryRunner func(queries []interface{}) ([]interface{}, error)

func (qb *QueryBuilder) Init(queryRunner QueryRunner, mixins ...*QueryMixin) {
	qb.queries = []interface{}{}
	qb.QueryRunner = queryRunner

	for _, mixin := range mixins {
		mixin.QueryBuilder = qb
	}
}

type QueryApplier func(entry interface{})

func (qb *QueryBuilder) Apply(applier QueryApplier) error {
	entries, err := qb.QueryRunner(qb.queries)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		applier(entry)
	}

	return nil
}
