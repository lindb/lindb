lexer grammar SQLLexer;

options { caseInsensitive = true; }

channels {
    COMMENT 
}

SIMPLE_COMMENT           : '--' ~[\r\n]* '\r'? '\n'? -> channel(COMMENT) ;
BRACKETED_COMMENT        : '/*' .*? '*/' -> channel(COMMENT) ;
WS                       : [ \r\n\t]+ -> channel(HIDDEN) ;

ALL                      : 'ALL' ;
ALIVE                    : 'ALIVE' ;
AND                      : 'AND' ;
ANALYZE                  : 'ANALYZE' ;
AS                       : 'AS' ;
ASC                      : 'ASC' ;
BETWEEN                  : 'BETWEEN' ;
BROKER                   : 'BROKER' ;
BROKERS                  : 'BROKERS' ;
BY                       : 'BY' ;
COMPACT                  : 'COMPACT' ;
CREATE                   : 'CREATE' ;
CROSS                    : 'CROSS' ;
COLUMNS                  : 'COLUMNS' ;
DATABASE                 : 'DATABASE' ;
DATABASES                : 'DATABASES' ;
DEFAULT                  : 'DEFAULT' ;
DESC                     : 'DESC' ;
DISTRIBUTED              : 'DISTRIBUTED' ;
DROP                     : 'DROP' ;
ENGINE                   : 'ENGINE' ;
ESCAPE                   : 'ESCAPE' ;
EXPLAIN                  : 'EXPLAIN' ;
EXISTS                   : 'EXISTS' ;
FALSE                    : 'FALSE' ;
FIELDS                   : 'FIELDS' ;
FLUSH                    : 'FLUSH' ;
FROM                     : 'FROM' ;
GROUP                    : 'GROUP' ;
HAVING                   : 'HAVING' ;
IF                       : 'IF' ;
IN                       : 'IN' ;
INTERVAL                 : 'INTERVAL' ;
JOIN                     : 'JOIN' ;
KEYS                     : 'KEYS' ;
LEFT                     : 'LEFT' ;
LIKE                     : 'LIKE' ;
LIMIT                    : 'LIMIT' ;
LOG                      : 'LOG' ;
LOGICAL                  : 'LOGICAL' ;
MASTER                   : 'MASTER' ;
METRICS                  : 'METRICS' ;
METRIC                   : 'METRIC' ;
METADATA                 : 'METADATA' ;
METADATAS                : 'METADATAS' ;
NAMESPACE                : 'NAMESPACE' ;
NAMESPACES               : 'NAMESPACES' ;
NOT                      : 'NOT' ;
NOW                      : 'NOW' ;
ON                       : 'ON' ;
OR                       : 'OR' ;
ORDER                    : 'ORDER' ;
REQUESTS                 : 'REQUESTS' ;
REPLICATIONS             : 'REPLICATIONS' ;
RIGHT                    : 'RIGHT' ;
ROLLUP                   : 'ROLLUP' ;
SELECT                   : 'SELECT' ;
SHOW                     : 'SHOW' ;
STATE                    : 'STATE' ;
STORAGE                  : 'STORAGE' ;
TABLE_NAMES              : 'TABLE_NAMES' ;
TIME                     : 'TIME' ;
TRACE                    : 'TRACE' ;
TRUE                     : 'TRUE' ;
TYPE                     : 'TYPE' ;
TYPES                    : 'TYPES' ;
VALUES                   : 'VALUES' ;
WHERE                    : 'WHERE' ;
WITH                     : 'WITH' ;
WITHIN                   : 'WITHIN' ;
USING                    : 'USING' ;
USE                      : 'USE' ;

// interval unit
SECOND                   : 'SECOND' ;
MINUTE                   : 'MINUTE' ;
HOUR                     : 'HOUR' ;
DAY                      : 'DAY' ;
MONTH                    : 'MONTH' ;
YEAR                     : 'YEAR' ;

EQ    : '=' ;
NEQ   : '<>' | '!=' ;
LT    : '<' ;
LTE   : '<=' ;
GT    : '>' ;
GTE   : '>=' ;

PLUS     : '+' ;
MINUS    : '-' ;
ASTERISK : '*' ;
SLASH    : '/' ;
PERCENT  : '%' ;

REGEXP             : '=~' ;
NEQREGEXP          : '!~' ;

EXCLAMATION_SYMBOL : '!' ;
DOT                : '.' ;
LR_BRACKET         : '(' ;
RR_BRACKET         : ')' ;
COMMA              : ',' ;

STRING                        : '\'' ( ~'\'' | '\'\'' )* '\'' ;
INTEGER_VALUE                 : DECIMAL_INTEGER ;
DECIMAL_VALUE                 : DECIMAL_INTEGER '.' DECIMAL_INTEGER?
                              | '.' DECIMAL_INTEGER
                              ;
DOUBLE_VALUE                  : DIGIT+ ('.' DIGIT*)? EXPONENT
                              | '.' DIGIT+ EXPONENT
                              ;
IDENTIFIER                    : (LETTER | '_') (LETTER | DIGIT | '_')* ;
DIGIT_IDENTIFIER              : DIGIT (LETTER | DIGIT | '_')+ ;
QUOTED_IDENTIFIER             : '"' ( ~'"' | '""' )* '"' ;
BACKQUOTED_IDENTIFIER         : '`' ( ~'`' | '``' )* '`';

// fragments for literal primitives
fragment DECIMAL_INTEGER      : DIGIT ('_'? DIGIT)* ;
fragment EXPONENT             : 'E' [+-]? DIGIT+ ;
fragment DIGIT                : [0-9] ;
fragment LETTER               : [A-Z] ;
