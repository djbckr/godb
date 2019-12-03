package dml

import (
	"fmt"
	"github.com/djbckr/godb/sql/token"
)

/*
https://docs.oracle.com/cd/E11882_01/server.112/e41084/statements_10002.htm

subquery [ for_update_clause ] ;

subquery ::=
{ query_block
| subquery { UNION [ALL] | INTERSECT | MINUS } subquery
    [ { UNION [ALL] | INTERSECT | MINUS } subquery ]...
| ( subquery )
} [ order_by_clause ]

query_block ::=
  [ subquery_factoring_clause ]
SELECT [ hint ] [ { { DISTINCT | UNIQUE } | ALL } ] select_list
  FROM { table_reference | join_clause | ( join_clause ) }
         [ , { table_reference | join_clause | (join_clause) } ] ...
  [ where_clause ]
  [ hierarchical_query_clause ]
  [ group_by_clause ]
  [ model_clause ]

!!!! NOTE !!!!
Support SELECT -> FROM or FROM -> SELECT
Also see Expressions (expr ::=) https://docs.oracle.com/cd/E11882_01/server.112/e41084/expressions001.htm
Also see Analytic Functions (the OVER keyword) https://docs.oracle.com/cd/E11882_01/server.112/e41084/functions004.htm

subquery_factoring_clause ::=
WITH
  query_name ([c_alias [, c_alias]...]) AS (subquery) [search__clause] [cycle_clause]
  [, query_name ([c_alias [, c_alias]...]) AS (subquery) [search_clause] [cycle_clause]]...

search__clause ::=
{ SEARCH
        { DEPTH FIRST BY c_alias [, c_alias]...
            [ ASC | DESC ]
            [ NULLS FIRST | NULLS LAST ]
         | BREADTH FIRST BY c_alias [, c_alias]...
            [ ASC | DESC ]
            [ NULLS FIRST | NULLS LAST ]
        }
        SET ordering_column
}

cycle_clause ::=
{CYCLE c_alias [, c_alias]...
    SET cycle_mark_c_alias TO cycle_value
    DEFAULT no_cycle_value
}

select_list ::=
{ [t_alias.] *
     | { query_name.* | [ schema. ] { table | view | materialized view } .* | expr [ [ AS ] c_alias ] }
    [, { query_name.* | [ schema. ] { table | view | materialized view } .* | expr [ [ AS ] c_alias ] } ]...
}

table_reference ::=
{ ONLY (query_table_expression)
| query_table_expression [ pivot_clause | unpivot_clause ]
} [ flashback_query_clause ]
  [ t_alias ]

flashback_query_clause ::=
{ VERSIONS BETWEEN
  { SCN | TIMESTAMP }
  { expr | MINVALUE } AND { expr | MAXVALUE }
| AS OF { SCN | TIMESTAMP } expr
}

query_table_expression ::=
{ query_name
| [ schema. ]
  { table [ partition_extension_clause | @ dblink ]
  | { view | materialized view } [ @ dblink ]
  } [sample_clause]
| (subquery [ subquery_restriction_clause ])
| table_collection_expression
}

subquery_restriction_clause :==
WITH { READ ONLY
     | CHECK OPTION
     } [ CONSTRAINT constraint ]

table_collection_expression ::=
TABLE (collection_expression) [ (+) ]

join_clause ::=
table_reference
  { inner_cross_join_clause | outer_join_clause }...

inner_cross_join_clause ::=
{ [ INNER ] JOIN table_reference
    { ON condition
    | USING (column [, column ]...)
    }
| { CROSS
  | NATURAL [ INNER ]
  }
  JOIN table_reference
}

outer_join_clause ::=
  [ query_partition_clause ] [ NATURAL ]
outer_join_type JOIN table_reference
  [ query_partition_clause ]
  [ ON condition
  | USING ( column [, column ]...)
  ]

query_partition_clause ::=
PARTITION BY
  { expr[, expr ]...
  | ( expr[, expr ]... )
  }

outer_join_type ::=
{ FULL | LEFT | RIGHT } [ OUTER ]

where_clause ::=
WHERE condition

condition ::=
https://docs.oracle.com/cd/E11882_01/server.112/e41084/conditions.htm#g1077361

group_by_clause ::=
GROUP BY
   { expr
   | rollup_cube_clause
   | grouping_sets_clause
   }
     [, { expr
        | rollup_cube_clause
        | grouping_sets_clause
        }
     ]...
   [ HAVING condition ]


order_by_clause ::=
ORDER [ SIBLINGS ] BY
{ expr | position | c_alias }
[ ASC | DESC ]
[ NULLS FIRST | NULLS LAST ]
  [, { expr | position | c_alias }
     [ ASC | DESC ]
     [ NULLS FIRST | NULLS LAST ]
  ]...

for_update_clause ::=
FOR UPDATE
  [ OF [ [ schema. ] { table | view } . ] column
         [, [ [ schema. ] { table | view } . ] column
         ]...
  ]
  [ { NOWAIT | WAIT integer
    |  SKIP LOCKED
    }
  ]

*/

type Query struct {
	QueryBlock []*TQueryBlock
}

type TSetOperator = int

const (
	NONE TSetOperator = iota
	UNION
	UNION_ALL
	INTERSECT
	MINUS
)

type TJoinType = int

const (
	JOIN_INNER TJoinType = iota
	JOIN_INNER_CROSS
	JOIN_INNER_NATURAL
	JOIN_OUTER_LEFT
	JOIN_OUTER_RIGHT
	JOIN_OUTER_FULL
)

type TQueryBlock struct {
	SetOperator TSetOperator
	From        []*TFrom
	Select      []*TSelect
	Where       []*TWhere
	GroupBy     []*TGroupBy
}

type TFrom struct {
	JoinType TJoinType
	TableRef *TTableRef
	On       []*TCondition
}

type TSelect struct {
}

type TCondition struct {
}

type TTableRef struct {
}

type TWhere struct {

}

type TGroupBy struct {
}

func ProcessSelect(sql token.Tokens) {
	var vv interface{}

	for idx, tkn := range sql {
		vv = tkn.Value
		fmt.Printf("%v, %t   ::   %v    ::    %v\n", idx, vv, vv, tkn.TokenType)
	}
}
