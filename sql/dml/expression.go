package dml
/*

expr ::=
{ simple_expression
| compound_expression
| case_expression
| cursor_expression
| datetime_expression
| function_expression
| interval_expression
| object_access_expression
| scalar_subquery_expression
| model_expression
| type_constructor_expression
| variable_expression
}

simple_expression ::=
{ [ query_name.
  | [schema.]
    { table. | view. | materialized view. }
  ] { column | ROWID }
| ROWNUM
| string
| number
| sequence. { CURRVAL | NEXTVAL }
| NULL
}


compound_expression ::=
{ (expr)
| { + | - | PRIOR } expr
| expr { * | / | + | - | || } expr
}

Note: The double vertical bars are part of the syntax
      (indicating concatenation) rather than BNF notation.

case_expression ::=
CASE { simple_case_expression
     | searched_case_expression
     }
     [ else_clause ]
     END

simple_case_expression ::=
expr
  { WHEN comparison_expr THEN return_expr }...

searched_case_expression ::=
{ WHEN condition THEN return_expr }...

else_clause ::=
ELSE else_expr

cursor_expression ::=
CURSOR (subquery)

datetime_expression ::=
expr AT
   { LOCAL
   | TIME ZONE { ' [ + | - ] hh:mi'
               | DBTIMEZONE
               | 'time_zone_name'
               | expr
               }
   }

function_expression ::=

interval_expression ::=
( expr1 - expr2 )
   { DAY [ (leading_field_precision) ] TO
     SECOND [ (fractional_second_precision) ]
   | YEAR [ (leading_field_precision) ] TO
     MONTH
   }

object_access_expression ::=
{ table_alias.column.
| object_table_alias.
| (expr).
}
{ attribute [.attribute ]...
  [.method ([ argument [, argument ]... ]) ]
| method ([ argument [, argument ]... ])
}

scalar_subquery_expression ::=
A scalar subquery expression is a subquery that returns exactly one column value from one row

model_expression ::=
// skip this for now

type_constructor_expression ::=
[ NEW ] [ schema. ]type_name
   ([ expr [, expr ]... ])

variable_expression ::=
{ expr [, expr ]...
| ( [expr [, expr ]] ...)
}


*/
