// Define a grammar called LinSQL for LinDB query language
grammar SQL;

statement                : statementList EOF;

statementList           : queryStmt;

//data query plan
queryStmt               : T_EXPLAIN? selectExpr fromClause whereClause? groupByClause? orderByClause? limitClause? T_WITH_VALUE?;
selectExpr              : T_SELECT fields;
//select fields
fields                   : field ( T_COMMA field )* ;
field                    : fieldExpr alias? ;
alias                    : T_AS ident ;

//from clause
fromClause              : T_FROM metricName ;

//where clause
whereClause             : T_WHERE conditionExpr;

conditionExpr          : tagFilterExpr | tagFilterExpr T_AND timeRangeExpr | timeRangeExpr (T_AND tagFilterExpr)?;

tagFilterExpr          :
                         T_OPEN_P tagFilterExpr T_CLOSE_P
                       | tagKey (T_EQUAL | T_LIKE | T_NOT T_LIKE | T_REGEXP | T_NEQREGEXP | T_NOTEQUAL | T_NOTEQUAL2) tagValue
                       | tagKey (T_IN | T_NOT T_IN) T_OPEN_P tagValueList T_CLOSE_P
                       | tagFilterExpr (T_AND | T_OR) tagFilterExpr
                       ;

tagValueList           : tagValue (T_COMMA tagValue)*;
timeRangeExpr          : timeExpr (T_AND timeExpr)? ;
timeExpr               : T_TIME binaryOperator (nowExpr | ident) ;

nowExpr                 : nowFunc  durationLit? ;

nowFunc                 : T_NOW T_OPEN_P exprFuncParams? T_CLOSE_P ;

//group by
groupByClause          : T_GROUP T_BY groupByKeys (T_FILL T_OPEN_P fillOption T_CLOSE_P)? havingClause? ;
groupByKeys            : groupByKey (T_COMMA groupByKey)* ;
groupByKey             : ident | T_TIME T_OPEN_P durationLit T_CLOSE_P ;
fillOption             : T_NULL | T_PREVIOUS | L_INT | L_DEC ;

orderByClause          : T_ORDER T_BY sortFields ;
sortField               : fieldExpr ( T_ASC | T_DESC )* ;
sortFields              : sortField ( T_COMMA sortField )* ;

havingClause            : T_HAVING boolExpr ;
boolExpr                :
                           T_OPEN_P boolExpr T_CLOSE_P
                         | boolExpr boolExprLogicalOp boolExpr
                         | boolExprAtom
                         ;
boolExprLogicalOp     : T_AND  | T_OR ;
boolExprAtom           : binaryExpr ;
binaryExpr         : fieldExpr binaryOperator fieldExpr;
binaryOperator:
                           T_EQUAL
                         | T_NOTEQUAL
                         | T_NOTEQUAL2
                         | T_LESS
                         | T_LESSEQUAL
                         | T_GREATER
                         | T_GREATEREQUAL
                         | (T_LIKE | T_REGEXP)
                         ;

fieldExpr                :
                           fieldExpr T_MUL fieldExpr
                         | fieldExpr T_DIV fieldExpr
                         | fieldExpr T_ADD fieldExpr
                         | fieldExpr T_SUB fieldExpr
                         | T_OPEN_P fieldExpr T_CLOSE_P
                         | exprFunc
                         | exprAtom
                         | durationLit
                         ;

durationLit             : intNumber intervalItem ;
intervalItem            :
                           T_SECOND
                         | T_MINUTE
                         | T_HOUR
                         | T_DAY
                         | T_WEEK
                         | T_MONTH
                         | T_YEAR
                         ;
exprFunc                : funcName T_OPEN_P exprFuncParams? T_CLOSE_P ;
funcName                : T_SUM | T_MIN | T_MAX | T_AVG | T_STDDEV | T_HISTOGRAM;
exprFuncParams          : funcParam (T_COMMA funcParam)* ;
funcParam               :
                           fieldExpr
                         | tagFilterExpr
                         ;
exprAtom                :
                           ident identFilter?
                         | decNumber
                         | intNumber
                         ;
identFilter             : T_OPEN_SB tagFilterExpr T_CLOSE_SB ;

// Integer (positive or negative)
intNumber               : ('-' | '+')? L_INT ;
// Decimal number (positive or negative)
decNumber               : ('-' | '+')? L_DEC ;
limitClause             : T_LIMIT L_INT ;
metricName              : ident ;
tagKey                  : ident ;
tagValue                : ident ;
ident                    :  (L_ID | nonReservedWords) ('.' (L_ID | nonReservedWords))* ;


nonReservedWords      :
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
                        | T_SUM
                        | T_MIN
                        | T_MAX
                        | T_AVG
                        | T_STDDEV
                        | T_HISTOGRAM
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

T_SUM                : S U M                            ;
T_MIN                : M I N                            ;
T_MAX                : M A X                            ;
T_AVG                : A V G                            ;
T_STDDEV             : S T D D E V                      ;
T_HISTOGRAM          : H I S T O G R A M                ;

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
T_NEQREGEXP          :  '!~'  ;
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