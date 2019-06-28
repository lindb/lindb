// Define a grammar called LinSQL for LinDB query language
grammar SQL;

statement                : statement_list EOF;

statement_list           :
                           create_database_stmt
                         | update_database_stmt
                         | drop_database_stmt
                         | show_databases_stmt
                         | show_node_stmt
                         | show_measurements_stmt
                         | show_info_stmt
                         | show_tag_keys_stmt
                         | show_queries_stmt
                         | show_tag_values_stmt
                         | show_tag_values_info_stmt
                         | show_field_keys_stmt
                         | show_stats_stmt
                         | kill_query_stmt
                         | query_stmt
                         ;

create_database_stmt     : T_CREATE T_DATASBAE database_name ( T_WITH with_clause_list )? (T_COMMA interval_define_list)?;
with_clause_list         : with_clause (T_COMMA with_clause)* ;
with_clause              :
                           (T_INTERVAL duration_lit)
                         | (T_SHARD shard_num)
                         | (T_REPLICATION replica_factor)
                         | (T_TTL ttl_val)
                         | (T_META_TTL metattl_val)
                         | (T_PAST_TTL past_val)
                         | (T_FUTURE_TTL future_val)

                         ;
interval_define_list     : interval_define (T_COMMA interval_define)* ;
interval_define          : T_OPEN_P T_INTERVAL_NAME interval_name_val T_COMMA T_TTL ttl_val T_COMMA T_INTERVAL duration_lit  T_CLOSE_P;
shard_num                : int_number;
ttl_val                  : duration_lit;
metattl_val              : duration_lit;
past_val                 : duration_lit;
future_val               : duration_lit;
interval_name_val        : ident;
replica_factor           : int_number;
database_name            : ident;

update_database_stmt     : T_UPDATE T_DATASBAE database_name ( T_WITH with_clause_list )? (T_COMMA interval_define_list)?;

drop_database_stmt       : T_DROP T_DATASBAE database_name;

//meta data query plan
show_databases_stmt      : T_SHOW T_DATASBAES;
show_node_stmt           : T_SHOW T_NODE;
show_measurements_stmt   : T_SHOW T_MEASUREMENTS with_measurement_clause? limit_clause? ;
show_tag_keys_stmt       : T_SHOW T_TAG T_KEYS T_FROM metric_name limit_clause? ;
show_info_stmt           : T_SHOW T_INFO T_FROM metric_name;
show_tag_values_stmt     : T_SHOW T_TAG T_VALUES T_FROM metric_name with_tag_clause where_tag_cascade?  limit_clause? ;
show_tag_values_info_stmt: T_SHOW T_TAG T_VALUES T_INFO T_FROM metric_name with_tag_clause where_tag_cascade;
show_field_keys_stmt     : T_SHOW T_FIELD T_KEYS T_FROM metric_name limit_clause? ;
show_queries_stmt        : T_SHOW T_QUERIES limit_clause? ;
show_stats_stmt          : T_SHOW T_STATS ( T_FOR module)? (T_WITH component)?;
with_measurement_clause  : T_WITH T_MEASUREMENT ( T_EQUAL metric_name | T_REGEXP metric_name ) ;
with_tag_clause          : T_WITH T_KEY T_EQUAL tag_key;
where_tag_cascade        : T_WHERE  tag_cascade_expr;
//kill query plan
kill_query_stmt          : T_KILL T_QUERY query_id (T_ON server_id)? ;
query_id                 : L_INT ;
server_id                : L_INT ;
module                   : ident ;
component                : ident ;

//data query plan
query_stmt               : T_EXPLAIN? T_SELECT fields from_clause where_clause? group_by_clause? interval_by_clause? order_by_clause? limit_clause? T_WITH_VALUE?;
//select fields
fields                   : field ( T_COMMA field )* ;
field                    : expr alias? ;
alias                    : T_AS ident ;

//from clause
from_clause              : T_FROM metric_name ;

//where clause
where_clause             : T_WHERE clause_boolean_expr;

clause_boolean_expr      :
                           tag_boolean_expr
                         | time_expr
                         | clause_boolean_expr T_AND clause_boolean_expr
                         ;
tag_cascade_expr         : tag_equal_expr
                         | tag_boolean_expr
                         | tag_equal_expr (T_AND  tag_boolean_expr)?
                         ;
tag_equal_expr           : T_VALUE T_EQUAL tag_value_pattern ;

tag_boolean_expr         :
                           T_OPEN_P tag_boolean_expr T_CLOSE_P
                         | tag_key (T_EQUAL | T_LIKE | T_REGEXP | T_NOTEQUAL | T_NOTEQUAL2) tag_value
                         | tag_key (T_IN | T_NOT T_IN) T_OPEN_P tag_value_list T_CLOSE_P
                         | tag_boolean_expr (T_AND | T_OR) tag_boolean_expr
                         ;
tag_value_list           : tag_value (T_COMMA tag_value)*;
time_expr                : time_boolean_expr (T_AND time_boolean_expr)? ;
time_boolean_expr        : T_TIME bool_expr_binary_operator (now_expr | ident) ;

now_expr                 : now_func  duration_lit? ;

now_func                 : T_NOW T_OPEN_P expr_func_params? T_CLOSE_P ;

//group by
group_by_clause          : T_GROUP T_BY dimensions (T_FILL T_OPEN_P fill_option T_CLOSE_P)? having_clause? ;
dimensions               : dimension (T_COMMA dimension)* ;
dimension                : ident | T_TIME T_OPEN_P duration_lit T_CLOSE_P ;
fill_option              : T_NULL | T_PREVIOUS | L_INT | L_DEC ;

order_by_clause          : T_ORDER T_BY sort_fields ;
interval_by_clause       : T_INTERVAL T_BY interval_name_val;
sort_field               : expr ( T_ASC | T_DESC )* ;
sort_fields              : sort_field ( T_COMMA sort_field )* ;

having_clause            : T_HAVING bool_expr ;
bool_expr                :
                           T_OPEN_P bool_expr T_CLOSE_P
                         | bool_expr bool_expr_logical_op bool_expr
                         | bool_expr_atom
                         ;
bool_expr_logical_op     : T_AND  | T_OR ;
bool_expr_atom           : bool_expr_binary ;
                         //bool_expr_single_in
bool_expr_binary         : expr bool_expr_binary_operator expr;
bool_expr_binary_operator:
                           T_EQUAL
                         | T_NOTEQUAL
                         | T_NOTEQUAL2
                         | T_LESS
                         | T_LESSEQUAL
                         | T_GREATER
                         | T_GREATEREQUAL
                         | (T_LIKE | T_REGEXP)
                         ;

expr                     :
                           expr T_MUL expr
                         | expr T_DIV expr
                         | expr T_ADD expr
                         | expr T_SUB expr
                         | T_OPEN_P expr T_CLOSE_P
                         | expr_func
                         | expr_atom
                         | duration_lit
                         ;

duration_lit             : int_number interval_item ;
interval_item            :
                           T_SECOND
                         | T_MINUTE
                         | T_HOUR
                         | T_DAY
                         | T_WEEK
                         | T_MONTH
                         | T_YEAR
                         ;
expr_func                : ident T_OPEN_P expr_func_params? T_CLOSE_P ;
expr_func_params         : func_param (T_COMMA func_param)* ;
func_param               :
                           expr
                         | tag_boolean_expr
                         ;
expr_atom                :
                           ident ident_filter?
                         | dec_number
                         | int_number
                         ;
ident_filter             :
                            T_OPEN_SB tag_boolean_expr T_CLOSE_SB ;
//ident_conditon           : T_OPEN_SB tag_boolean_expr T_CLOSE_SB ;

// Integer (positive or negative)
int_number               : ('-' | '+')? L_INT ;
// Decimal number (positive or negative)
dec_number               : ('-' | '+')? L_DEC ;
limit_clause             : T_LIMIT L_INT ;
metric_name              : ident ;
tag_key                  : ident ;
tag_value                : ident ;
tag_value_pattern        : ident ;
ident                    :  (L_ID | non_reserved_words) ('.' (L_ID | non_reserved_words))* ;


non_reserved_words      :
                          T_CREATE
                        | T_INTERVAL
                        | T_SHARD
                        | T_REPLICATION
                        | T_TTL
                        | T_DATASBAE
                        | T_KILL
                        | T_SHOW
                        | T_DATASBAES
                        | T_NODE
                        | T_MEASUREMENTS
                        | T_MEASUREMENT
                        | T_FIELD
                        | T_TAG
                        | T_KEYS
                        | T_KEY
                        | T_WITH
                        | T_VALUES
                        | T_FROM
                        | T_WHERE
                        | T_LIMIT
                        | T_QUERIES
                        | T_QUERY
                        | T_SELECT
                        | T_AS
                        | T_AND
                        | T_OR
                        | T_NULL
                        | T_PREVIOUS
                        | T_FILL
                        | T_ORDER
                        | T_ASC
                        | T_DESC
                        | T_LIKE
                        | T_NOT
                        | T_BETWEEN
                        | T_IS
                        | T_PROFILE
                        | T_GROUP
                        | T_BY
                        | T_ON
                        | T_STATS
                        | T_TIME
                        | T_FOR
                        | T_SECOND
                        | T_MINUTE
                        | T_HOUR
                        | T_DAY
                        | T_WEEK
                        | T_MONTH
                        | T_YEAR
                        ;

// Lexer rules
T_CREATE             : C R E A T E                      ;
T_UPDATE             : U P D A T E                      ;
T_SET                : S E T                            ;
T_DROP               : D R O P                          ;
T_INTERVAL           : I N T E R V A L                  ;
T_INTERVAL_NAME      : N A M E                          ;
T_SHARD              : S H A R D                        ;
T_REPLICATION        : R E P L I C A T I O N            ;
T_TTL                : T T L                            ;
T_META_TTL           : M E T A T T L                    ;
T_PAST_TTL           : P A S T T T L                    ;
T_FUTURE_TTL         : F U T U R E T T L                ;
T_KILL               : K I L L                          ;
T_ON                 : O N                              ;
T_SHOW               : S H O W                          ;
T_DATASBAE           : D A T A B A S E                  ;
T_DATASBAES          : D A T A B A S E S                ;
T_NODE               : N O D E                          ;
T_MEASUREMENTS       : M E A S U R E M E N T S          ;
T_MEASUREMENT        : M E A S U R E M E N T            ;
T_FIELD              : F I E L D                        ;
T_TAG                : T A G                            ;
T_INFO               : I N F O                          ;
T_KEYS               : K E Y S                          ;
T_KEY                : K E Y                            ;
T_WITH               : W I T H                          ;
T_VALUES             : V A L U E S                      ;
T_VALUE              : V A L U E                        ;
T_FROM               : F R O M                          ;
T_WHERE              : W H E R E                        ;
T_LIMIT              : L I M I T                        ;
T_QUERIES            : Q U E R I E S                    ;
T_QUERY              : Q U E R Y                        ;
T_EXPLAIN            : E X P L A I N                    ;
T_WITH_VALUE         : W I T H V A L U E                ;
T_SELECT             : S E L E C T                      ;
T_AS                 : A S                              ;
T_AND                : A N D                            ;
T_OR                 : O R                              ;
T_FILL               : F I L L                          ;
T_NULL               : N U L L                          ;
T_PREVIOUS           : P R E V I O U S                  ;
T_ORDER              : O R D E R                        ;
T_ASC                : A S C                            ;
T_DESC               : D E S C                          ;
T_LIKE               : L I K E                          ;
T_NOT                : N O T                            ;
T_BETWEEN            : B E T W E E N                    ;
T_IS                 : I S                              ;
T_GROUP              : G R O U P                        ;
T_HAVING             : H A V I N G                      ;
T_BY                 : B Y                              ;
T_FOR                : F O R                            ;
T_STATS              : S T A T S                        ;
T_TIME               : T I M E                          ;
T_NOW                : N O W                            ;
T_IN                 : I N                              ;

T_LOG                : L O G                            ;
T_PROFILE            : P R O F I L E                    ;


//time unit
T_SECOND             : S                                ;
T_MINUTE             : 'm'                              ;
T_HOUR               : H                                ;
T_DAY                : D                                ;
T_WEEK               : W                                ;
T_MONTH              : 'M'                              ;
T_YEAR               : Y                                ;

//
T_DOT                :  '.'   ;
T_COLON              :  ':'   ;
T_EQUAL              :  '='   ;
T_NOTEQUAL           :  '<>'  ;
T_NOTEQUAL2          :  '!='  ;
T_GREATER            :  '>'   ;
T_GREATEREQUAL       :  '>='  ;
T_LESS               :  '<'   ;
T_LESSEQUAL          :  '<='  ;
T_REGEXP             :  '=~'  ;
T_COMMA              :  ','   ;
T_OPEN_B             :  '{'   ;
T_CLOSE_B            :  '}'   ;
T_OPEN_SB            :  '['   ;
T_CLOSE_SB           :  ']'   ;
T_OPEN_P             :  '('   ;
T_CLOSE_P            :  ')'   ;
T_ADD                :  '+'   ;
T_SUB                :  '-'   ;
T_DIV                :  '/'   ;
T_MUL                :  '*'   ;
T_MOD                :  '%'   ;

L_ID                 : L_ID_PART ;
L_INT                : L_DIGIT+       ;                                               // Integer
L_DEC                : L_DIGIT+ '.' ~'.' L_DIGIT*                               // Decimal number
                     | '.' L_DIGIT+
                     ;
WS                   : BLANK+ -> skip ;                                        //Whitespace

fragment BLANK       : [ \t\r\n]      ;
fragment L_DIGIT     : [0-9] ;
// Double quoted string escape sequence
//fragment L_STR_ESC_D : '""' | '\\"' ;
fragment L_ID_PART   :
                      [a-zA-Z] ([a-zA-Z] | L_DIGIT | '_' | '.')*                                            // Identifier part
                      | '$' '{' .*? '}'
                      | ('_' | '@' | ':' | '#' | '$') ([a-zA-Z] | L_DIGIT | '_' | '@' | ':' | '#' | '$')+     // (at least one char must follow special char)
                      | '"' .*? '"'                                                                           // Quoted identifiers
                      | '`' .*? '`'                                                                           // Quoted identifiers
                      | '\'' .*? '\''                                                                           // Quoted identifiers
                     ;

// Support case-insensitive keywords and allowing case-sensitive identifiers
fragment A : ('a'|'A') ;
fragment B : ('b'|'B') ;
fragment C : ('c'|'C') ;
fragment D : ('d'|'D') ;
fragment E : ('e'|'E') ;
fragment F : ('f'|'F') ;
fragment G : ('g'|'G') ;
fragment H : ('h'|'H') ;
fragment I : ('i'|'I') ;
fragment J : ('j'|'J') ;
fragment K : ('k'|'K') ;
fragment L : ('l'|'L') ;
fragment M : ('m'|'M') ;
fragment N : ('n'|'N') ;
fragment O : ('o'|'O') ;
fragment P : ('p'|'P') ;
fragment Q : ('q'|'Q') ;
fragment R : ('r'|'R') ;
fragment S : ('s'|'S') ;
fragment T : ('t'|'T') ;
fragment U : ('u'|'U') ;
fragment V : ('v'|'V') ;
fragment W : ('w'|'W') ;
fragment X : ('x'|'X') ;
fragment Y : ('y'|'Y') ;
fragment Z : ('z'|'Z') ;