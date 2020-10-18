package trigram_test

import (
	"testing"

	. "github.com/brownchow/trigram"
)

func TestTrigramlize(t *testing.T) {
	// 把string转成trigram数组, 字符串长度是n，最终的trigram数组的长度是n-2
	ret := ExtractStringToTrigram("Cod")
	t.Logf("trigram of Cod is: %d ", ret)
	if ret[0] != 4419428 {
		t.Errorf("Trigram failed, expect 4419428\n")
	}
	//string length longer than 3
	ret = ExtractStringToTrigram("Code")
	// "Code"转成trigram之后是：4419428，7300197
	if ret[0] != 4419428 && ret[1] != 7300197 {
		t.Errorf("Trigram failed on longer string")
	}
}

func TestMapIntersect(t *testing.T) {
	mapA := make(map[int]bool)
	mapB := make(map[int]bool)

	mapA[1] = true
	mapA[2] = true
	mapB[1] = true

	ret := IntersectTwoMap(mapA, mapB)
	t.Logf("ret res is: %v", ret)
	if len(ret) != 1 || ret[1] == false {
		t.Errorf("Map intersect error")
	}

	ret = IntersectTwoMap(mapB, mapB)
	t.Logf("ret res is: %v", ret)
	if len(ret) != 1 || ret[1] == false {
		t.Errorf("Map intersect error")
	}

	mapA[3] = true
	mapB[3] = true
	mapA[4] = true

	ret = IntersectTwoMap(mapB, mapA)
	t.Logf("ret res is: %v", ret)
	if len(ret) != 2 || ret[1] == false {
		t.Errorf("Map intersect error")
	}
}

func TestTrigramIndexBasicQuery(t *testing.T) {
	ti := NewTrigramIndex()
	ti.Add("Code is my life")
	ti.Add("Search")
	ti.Add("I write a lot of Codes")

	ret := ti.Query("Code")
	t.Logf("search result is: %v", ret)
	if ret[0] != 1 || ret[1] != 3 {
		t.Errorf("Basic query is failed")
	}
}

func TestEmptyLessQuery(t *testing.T) {
	ti := NewTrigramIndex()
	ti.Add("Code is my life")
	ti.Add("Search")
	ti.Add("I write a lot of Codes")

	// 长度小于3的查询，会返回所有的文档
	ret := ti.Query("te") // less than 3, should get all doc ID
	if len(ret) != 3 || ret[0] != 1 || ret[2] != 3 {
		t.Errorf("Error on less than 3 character query")
	}

	ret = ti.Query("")
	if len(ret) != 3 || ret[0] != 1 || ret[2] != 3 {
		t.Errorf("Error on empty character query")
	}
}

func TestDelete(t *testing.T) {
	ti := NewTrigramIndex()
	ti.Add("Code is my life")
	// 删除所有跟含有 Code 的文档
	ti.Delete("Code", 1)
	ret := ti.Query("Code")
	if len(ret) != 0 {
		t.Error("Basic delete failed", ret)
	}
	// 因为前面已经把包含Code的文档删除了，所以这里就查询不到了
	ret = ti.Query("life")
	if len(ret) != 1 || ret[0] != 1 {
		t.Error("Basic delete failed", ret)
	}
}

func BenchmarkDelete(b *testing.B) {
	big := NewTrigramIndex()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Add("1234567890")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Delete("1234567890", i)
	}
}

func BenchmarkQuery(b *testing.B) {
	big := NewTrigramIndex()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Add("1234567890")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Query("1234567890")
	}
}

func BenchmarkIntersection(b *testing.B) {
	DocA := make(map[int]bool)
	DocB := make(map[int]bool)
	for i := 0; i < 101; i++ {
		DocA[i] = true
		DocB[i+1] = true
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntersectTwoMap(DocA, DocB)
	}
}
