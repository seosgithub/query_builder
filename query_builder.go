package query_builder

type IQueryMixin interface {
	GetQueryMixin() *QueryMixin
}

type QueryMixin struct {
	*QueryBuilder

	IQueryMixin
}

func (q *QueryMixin) GetQueryMixin() *QueryMixin {
	return q
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

func (qb *QueryBuilder) Init(queryRunner QueryRunner, mixins ...IQueryMixin) {
	qb.queries = []interface{}{}
	qb.QueryRunner = queryRunner

	for _, mixin := range mixins {
		mixin.GetQueryMixin().QueryBuilder = qb
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
