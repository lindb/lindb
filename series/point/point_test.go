package point_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/point"
	"github.com/lindb/lindb/series/tag"

	"github.com/stretchr/testify/assert"
)

func Test_Point_Tag_Name(t *testing.T) {
	var p point.Point
	assert.False(t, p.HasTag([]byte("not-exist")))
	p.SetName("cpu.load")
	assert.Equal(t, "cpu.load", string(p.Name()))
	assert.Equal(t, uint64(17143933640780095220), p.HashID())
	assert.Equal(t, uint64(17143933640780095220), p.HashID())
	p.SetName("memory")
	assert.Equal(t, uint64(2901243900060439616), p.HashID())
	assert.Equal(t, "memory", string(p.Name()))
	assert.Equal(t, "memory", string(p.Key()))
	assert.Equal(t, "", string(p.TagsHashKey()))

	p.AddTag("ip", "1. 1.1.1")
	assert.Equal(t, uint64(0xa5650c48f713ed35), p.TagsHashID())
	assert.Equal(t, uint64(0xa5650c48f713ed35), p.TagsHashID())
	p.AddTags("zone", "sh", "host", "t,est")
	p.AddTags("drop") //this tag won't be added
	assert.Equal(t, uint64(0xb43e3500622cf8f6), p.TagsHashID())
	assert.Equal(t, "memory,host=t\\,est,ip=1.\\ 1.1.1,zone=sh", string(p.Key()))
	assert.Len(t, p.Tags(), 3)
	assert.Equal(t, tag.Tags{
		{Key: []byte("host"), Value: []byte("t,est")},
		{Key: []byte("ip"), Value: []byte("1. 1.1.1")},
		{Key: []byte("zone"), Value: []byte("sh")},
	}, p.Tags())
	var counter int
	p.ForEachTag(func(k, v []byte) bool {
		if bytes.Equal(k, []byte("zone")) {
			return false
		}
		counter++
		return true
	})
	assert.Equal(t, 2, counter)

	assert.Equal(t, ",host=t\\,est,ip=1.\\ 1.1.1,zone=sh", string(p.TagsHashKey()))
	assert.True(t, p.HasTag([]byte("host")))
	assert.False(t, p.HasTag([]byte("not-exist")))
	p.SetTags(nil)
	assert.Nil(t, p.TagsHashKey())
	assert.Equal(t, "memory", string(p.Name()))
}

func Test_Point_Time(t *testing.T) {
	var p point.Point
	var buf []byte
	p.SetName("memory").AddTags("host", "1.1.1.1")
	assert.Equal(t, "memory,host=1.1.1.1 ", p.String())
	assert.Len(t, p.AppendString(buf), 20)
	assert.Equal(t, 20, p.StringSize())

	assert.Zero(t, p.UnixMilli())
	p.SetUnixMilli(1576904537)
	t1 := p.Time()
	assert.Equal(t, int64(1576904537000), p.UnixMilli())
	p.SetUnixMilli(1576904537000)
	assert.Equal(t, "memory,host=1.1.1.1 1576904537000", p.String())
	assert.Equal(t, 34, p.StringSize())
	buf = nil
	assert.Len(t, p.AppendString(buf), 33)

	t2 := p.Time()
	assert.Equal(t, t1, t2)
	p.SetTime(time.Now())
	assert.NotEqual(t, t2, p.UnixMilli())

	// negative time
	p.SetUnixMilli(-1000)
	assert.Equal(t, 29, p.StringSize())
}

func Test_Point_Fields(t *testing.T) {
	var p point.Point
	p.SetName("memory").
		AddTags("host", "1.1.1.1").
		AddField("gc Num", field.SumField, 100).
		AddField("timerSum", field.HistogramField, 10201).
		SetUnixMilli(100000000)
	assert.Equal(t,
		"memory,host=1.1.1.1 gc\\ Num_SUM=1E+02,timerSum_HGM=1.0201E+04 100000000000",
		p.String())
	_, err := p.Fields()
	assert.Nil(t, err)

	itr := p.FieldIterator()
	for itr.Next() {
		assert.Equal(t, "gc Num", string(itr.Name()))
		assert.Equal(t, field.SumField, itr.Type())
		break
	}

	p.SetFields(nil)
	itr = p.FieldIterator()
	assert.False(t, itr.Next())

	p.Reset()
}

func Test_Point_parseFields(t *testing.T) {
	p0 := new(point.Point).
		SetName("cpu .load").
		AddTags("ho st", "test", "i,p", "1.1.1.1", "zo=ne", "sh").
		AddField("last1min", field.SumField, 5).
		AddField("last5min", field.SumField, 10).
		AddField("last15min", field.SumField, 7).
		SetUnixMilli(1577000000)
	parsedPoints, _ := point.ParsePointsFromString(p0.String())
	parsedPoint := parsedPoints[0]
	fs, err := parsedPoint.Fields()
	assert.Nil(t, err)
	assert.Equal(t, field.SumField, fs[0].Type)
}

func Test_ParsePoints(t *testing.T) {
	p0 := new(point.Point).
		SetName("cpu .load").
		AddTags("ho st", "test", "i,p", "1.1.1.1", "zo=ne", "sh").
		AddField("last1min", field.SumField, 5).
		AddField("last5min", field.SumField, 10).
		AddField("last15min", field.SumField, 7).
		SetUnixMilli(1577000000)
	p1 := new(point.Point).
		SetName("memory").
		AddTags("host", "test2", "ip", "1.1.1.2", "zone", "bj").
		AddField("dirty", field.SumField, 10).
		AddField("cached", field.SumField, 30).
		AddField("total", field.SumField, 100).
		SetUnixMilli(1577000001)
	p2 := new(point.Point).
		SetName("tcp").
		AddTags("host", "test3", "ip", "1.1.1.3", "zone", "nj").
		AddField("established", field.SumField, -1).
		AddField("timeWaited", field.SumField, 2).
		AddField("closeWaited", field.SumField, 3).
		SetUnixMilli(1577000002)
	var buf []byte
	buf = p0.AppendString(buf)
	buf = append(buf, []byte("\n")[0])
	buf = p1.AppendString(buf)
	buf = append(buf, []byte("\n")[0])
	buf = p2.AppendString(buf)
	buf = append(buf, []byte("\n")[0])
	buf = append(buf, []byte("#test comment")...)

	points, err := point.ParsePointsFromString(string(buf))
	assert.Nil(t, err)
	assert.Equal(t, p0.String(), points[0].String())
	assert.Equal(t, p1.String(), points[1].String())
	assert.Equal(t, p2.String(), points[2].String())
}

func Test_ParsePoints_error(t *testing.T) {
	// metric name invalid
	_, err := point.ParsePointsFromString(
		"cpu. load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1577000000000")
	assert.NotNil(t, err)

	// escape metric-name
	_, err = point.ParsePointsFromString(
		"cpu.\\ load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1577000000000")
	assert.Nil(t, err)
	_, err = point.ParsePointsFromString(
		"cpu.\\=load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1577000000000")
	assert.Nil(t, err)

	// empty text
	_, err = point.ParsePointsFromString("")
	assert.Nil(t, err)

	// missing metric name
	_, err = point.ParsePointsFromString(
		",host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1577000000000")
	assert.NotNil(t, err)
}

func Test_ParsePoints_TagsError(t *testing.T) {
	// duplicate tags
	_, err := point.ParsePointsFromString(
		"cpu.load,host=test,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00")
	assert.NotNil(t, err)

	// sort tags
	points, err := point.ParsePointsFromString(
		"cpu.load,host=test,a=b,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1")
	assert.Nil(t, err)
	assert.Equal(t, "cpu.load,a=b,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1000",
		points[0].String())

	// invalid tag format
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=s=h last15min_SUM=7E+00 1577000000000")
	assert.NotNil(t, err)

	// tag value missing
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne= last15min_SUM=7E+00 1577000000000")
	assert.NotNil(t, err)
}

func Test_ParsePoints_FieldsError(t *testing.T) {
	// invalid number
	_, err := point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+i00 1577000000000")
	assert.NotNil(t, err)

	// escaped field-name
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15mi\\ n_xxx=7E+00 1577000000000")
	assert.Nil(t, err)

	// missing fields
	points, err := point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_xxx=7E+00 1577000000000")
	assert.Nil(t, err)
	_, err = points[0].Fields()
	assert.NotNil(t, err)

	// invalid number
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00xx")
	assert.NotNil(t, err)

	// invalid number
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=-")
	assert.NotNil(t, err)

	// NaN
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=NaN")
	assert.NotNil(t, err)

	// 2 points
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=1.1.1")
	assert.NotNil(t, err)
}

func Test_ParsePoints_TimestampError(t *testing.T) {
	// empty time
	_, err := point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 ")
	assert.Nil(t, err)
	points, err := point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00")
	assert.Nil(t, err)
	assert.NotZero(t, points[0].UnixMilli())

	// bad timestamp
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1.324")
	assert.NotNil(t, err)

	// scan to space
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1 ")
	assert.Nil(t, err)
	// scan to \n
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 1\n")
	assert.Nil(t, err)
	// negative timestamp
	_, err = point.ParsePointsFromString(
		"cpu.load,host=test,ip=1.1.1.1,zo\\=ne=sh last15min_SUM=7E+00 -1")
	assert.Nil(t, err)
}
