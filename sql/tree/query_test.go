package tree

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	stmt, err := GetParser().CreateStatement("select *,dd.*,s+dd/cc,sum(avg(a)+b,tt) as sum from (select sub as sub from table t group by all) sub_table group by (a,b,c) having sum>0 limit 10", NewNodeIDAllocator())
	assert.Equal(t, nil, err)
	b, err := json.MarshalIndent(stmt, "", "  ")
	assert.Equal(t, nil, err)
	fmt.Println(string(b))
}

func TestSelectStatement_With(t *testing.T) {
	stmt, err := GetParser().CreateStatement(`
with 
 a as (select a as aa from a where a in(a,c,b) group by dd),
 b as (select sum(b) from b where b=d) 
select 
	a.*,b.*,(a.sum+b.sum)/100
from 
	a left join b using(uid) 
	right join (select load from cpu) c on c=d 
where 
	a=1 and b like 'dd%' and (a=c or a=d) 
having sum>100 and total>10
order by a desc,b asc
limit 10
	`, NewNodeIDAllocator())
	assert.Equal(t, nil, err)
	b, err := json.MarshalIndent(stmt, "", "  ")
	assert.Equal(t, nil, err)
	fmt.Println(string(b))
}
