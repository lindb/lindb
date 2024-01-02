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
AS                       : 'AS' ;
ASC                      : 'ASC' ;
BROKER                   : 'BROKER' ;
BROKERS                  : 'BROKERS' ;
BY                       : 'BY' ;
COMPACT                  : 'COMPACT' ;
CREATE                   : 'CREATE' ;
CROSS                    : 'CROSS' ;
DATABASE                 : 'DATABASE' ;
DATABASES                : 'DATABASES' ;
DEFAULT                  : 'DEFAULT' ;
DESC                     : 'DESC' ;
DROP                     : 'DROP' ;
ESCAPE                   : 'ESCAPE' ;
EXPLAIN                  : 'EXPLAIN' ;
EXISTS                   : 'EXISTS' ;
FALSE                    : 'FALSE' ;
FIELDS                   : 'FIELDS' ;
FILTER                   : 'FILTER' ;
FLUSH                    : 'FLUSH' ;
FROM                     : 'FROM' ;
GROUP                    : 'GROUP' ;
HAVING                   : 'HAVING' ;
IF                       : 'IF' ;
IN                       : 'IN' ;
JOIN                     : 'JOIN' ;
KEYS                     : 'KEYS' ;
LEFT                     : 'LEFT' ;
LIKE                     : 'LIKE' ;
LIMIT                    : 'LIMIT' ;
MASTER                   : 'MASTER' ;
METRICS                  : 'METRICS' ;
METADATA                 : 'METADATA' ;
METADATAS                : 'METADATAS' ;
NAMESPACE                : 'NAMESPACE' ;
NAMESPACES               : 'NAMESPACES' ;
NOT                      : 'NOT' ;
ON                       : 'ON' ;
OR                       : 'OR' ;
ORDER                    : 'ORDER' ;
PLAN                     : 'PLAN' ;
REQUESTS                 : 'REQUESTS' ;
REPLICATIONS             : 'REPLICATIONS' ;
RIGHT                    : 'RIGHT' ;
ROLLUP                   : 'ROLLUP' ;
SELECT                   : 'SELECT' ;
SHOW                     : 'SHOW' ;
STATE                    : 'STATE' ;
STORAGE                  : 'STORAGE' ;
TAG                      : 'TAG' ;
TRUE                     : 'TRUE' ;
TYPES                    : 'TYPES' ;
VALUES                   : 'VALUES' ;
WHERE                    : 'WHERE' ;
WITH                     : 'WITH' ;
WITHIN                   : 'WITHIN' ;
USING                    : 'USING' ;
USE                      : 'USE' ;

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
