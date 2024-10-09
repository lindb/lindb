parser grammar SQLParser;

options {
    tokenVocab = SQLLexer;
}
// debug:(-gui/-tree)
// antlr4-parse SQLParser.g4 SQLLexer.g4 statement -gui
// select 10*(idle+10*5)/10,node from cpu group by node
// ^D
statement           : ddlStatement
                    | dmlStatement
                    | adminStatement
                    | utilityStatement 
                    EOF
                    ;

ddlStatement        : createDatabase
                    | dropDatabase
                    | createBroker 
                    ;

dmlStatement        : query                                                  #statementDefault
  									| EXPLAIN ('(' explainOption (',' explainOption)* ')')? 
                        dmlStatement                                         #explain
   				 					| EXPLAIN ANALYZE dmlStatement                           #explainAnalyze
                    ;

adminStatement      : flushDatabase
                    | compactDatabase 
                    | showStatement 
                    ;

utilityStatement    : useStatement ;

explainOption       : TYPE value=(LOGICAL | DISTRIBUTED)                     #explainType
										;                   

// ddl
createDatabase      : CREATE DATABASE name=qualifiedName
                       (WITH properties)?
                       (ROLLUP '(' rollupOptions (',' rollupOptions)* ')')? 
                    ;
rollupOptions       : properties ;
dropDatabase        : DROP DATABASE (IF EXISTS)? database=qualifiedName ;
createBroker        : CREATE BROKER name=qualifiedName
                       (WITH properties)?
                    ;

// administration
flushDatabase       : FLUSH DATABASE database=qualifiedName ;
compactDatabase     : COMPACT DATABASE database=qualifiedName ;
showStatement       : SHOW MASTER                                            #showMaster
                    | SHOW BROKERS                                           #showBrokers
                    | SHOW REQUESTS                                          #showRequests
                    | SHOW LIMIT                                             #showLimit
                    | SHOW METADATA TYPES                                    #showMetadataTypes
                    | SHOW METADATAS                                         #showMetadatas
                    | SHOW ALIVE                                             #showAlive
                    | SHOW REPLICATIONS                                      #showReplications
                    | SHOW STATE                                             #showState
                    | SHOW DATABASES                                         #showDatabases
                    ;

// utility
useStatement        : USE database=identifier ;

// dml
showMetricMetadata  : SHOW NAMESPACES                                        #showNamespaces
                    | SHOW METRICS                                           #showMetrics
                    | SHOW FIELDS                                            #showFields
                    | SHOW TAG KEYS                                          #showTagKeys
                    | SHOW TAG VALUES                                        #showTagValues
                    ;
query               : with? queryNoWith ;
with                : WITH namedQuery (',' namedQuery)* ;
namedQuery          : name=identifier AS '(' query ')' ;

//Removing redundant ORDER BY
//https://trino.io/blog/2019/06/03/redundant-order-by.html 
queryNoWith         : queryTerm  
                       (ORDER BY orderBy)?
                       (LIMIT limit=limitRowCount)?                          
                    ;

queryTerm           : queryPrimary                                           #queryTermDefault
                    ;

queryPrimary        : querySpecification                                     #queryPrimaryDefault
                    | '(' queryNoWith ')'                                    #subquery
                    ;

querySpecification  : SELECT selectItem (',' selectItem)*
                       (FROM relation (',' relation)*)?
                       (WHERE where=booleanExpression)?
                       (GROUP BY groupBy)?
                       (HAVING having)?
                    ;

selectItem          : expression (AS? identifier)?                           #selectSingle
                    | primaryExpression '.' ASTERISK                         #selectAll
                    | ASTERISK                                               #selectAll
                    ;
relation            : left=relation
                      ( CROSS JOIN right=relation
                        | joinType JOIN rightRelation=relation joinCriteria
                      )                                                     #joinRelation
                    | aliasedRelation                                       #relationDefault
                    ;

joinType            : LEFT | RIGHT ;
joinCriteria        : ON booleanExpression                                  
                    | USING '(' identifier (',' identifier)* ')'            
                    ;

aliasedRelation     : relationPrimary (AS? identifier)? ;
relationPrimary     : qualifiedName                                         #tableName
                    | '(' query ')'                                         #subQueryRelation
                    ;

groupBy             : groupingElement (',' groupingElement)* ;
groupingElement     : groupingSet                                           #singleGroupingSet
                    | ALL                                                   #groupByAllColumns
                    ;

groupingSet         : '(' (expression (',' expression)*)? ')'
                    | expression 
                    ;

having              : booleanExpression ;
orderBy             : sortItem (',' sortItem)* ;
sortItem            : expression ordering=(ASC | DESC)? ;

limitRowCount       : INTEGER_VALUE ; 

expression          : booleanExpression
                    ;
booleanExpression   : notOperator = (NOT | '!') booleanExpression           #logicalNot
                    | booleanExpression AND booleanExpression               #and
                    | booleanExpression OR booleanExpression                #or
                    | predicate                                             #predicatedExpression
                    ;
valueExpression     : primaryExpression                                     #valueExpressionDefault
                    | left=valueExpression operator=(ASTERISK | SLASH | PERCENT) right=valueExpression  #arithmeticBinary
                    | left=valueExpression operator=(PLUS | MINUS) right=valueExpression                #arithmeticBinary
                    ;
primaryExpression   : 
                    number                                                  #numericLiteral
                    | booleanValue                                          #booleanLiteral
                    | string                                                #stringLiteral
                    | qualifiedName '(' (expression (',' expression)*)? ')' #functionCall
                    | identifier                                            #columnReference
                    | base=primaryExpression '.' fieldName=identifier       #dereference
                    | '(' expression ')'                                    #parenExpression
                    ;

predicate           : left=valueExpression operator=comparisonOperator right=valueExpression                   #binaryComparisonPredicate
                    | left=valueExpression NOT? IN '(' expression (',' expression)* ')'                        #inPredicate
                    | left=valueExpression NOT? LIKE pattern=valueExpression (ESCAPE escape=valueExpression)?  #likePredicate
                    | left=valueExpression operator=(REGEXP|NEQREGEXP) pattern=valueExpression?                #regexpPredicate
                    | valueExpression                                                                          #valueExpressionPredicate
                    ;

comparisonOperator  : EQ | NEQ | LT | LTE | GT | GTE ;
filter              : FILTER '(' WHERE booleanExpression ')' ;

qualifiedName       : identifier ('.' identifier)* ;

properties          : '(' propertyAssignments ')' ;
propertyAssignments : property (',' property)* ;
property            : name=identifier EQ value=propertyValue ;
propertyValue       : DEFAULT                                               #defaultPropertyValue
                    | expression                                            #nonDefaultPropertyValue
                    ;

booleanValue        : TRUE | FALSE ;
string              : STRING                                                #basicStringLiteral
                    ;

identifier          : IDENTIFIER                                            #unquotedIdentifier
                    | QUOTED_IDENTIFIER                                     #quotedIdentifier
                    | nonReserved                                           #unquotedIdentifier
                    | BACKQUOTED_IDENTIFIER                                 #backQuotedIdentifier
                    | DIGIT_IDENTIFIER                                      #digitIdentifier
                    ;

number              : MINUS? DECIMAL_VALUE                                  #decimalLiteral
                    | MINUS? DOUBLE_VALUE                                   #doubleLiteral
                    | MINUS? INTEGER_VALUE                                  #integerLiteral
                    ;

nonReserved         :
                      ALL | ALIVE | AND | AS | ASC
                    | BROKER | BROKERS | BY 
                    | COMPACT | CREATE | CROSS 
                    | DATABASE | DATABASES | DEFAULT | DESC | DISTRIBUTED | DROP
                    | ESCAPE | EXPLAIN | EXISTS
                    | FALSE | FIELDS | FILTER | FLUSH | FROM
                    | GROUP 
                    | HAVING
                    | IF | IN 
                    | JOIN
                    | KEYS
                    | LEFT | LIKE | LIMIT | LOGICAL
                    | MASTER | METRICS | METADATA | METADATAS
                    | NAMESPACE | NAMESPACES | NOT
                    | ON | OR | ORDER
                    | PLAN
                    | REQUESTS | REPLICATIONS | RIGHT | ROLLUP
                    | SELECT | SHOW | STATE | STORAGE
                    | TAG | TRUE | TYPE | TYPES 
                    | VALUES
                    | WHERE | WITH | WITHIN
                    | USING | USE
                    ;

