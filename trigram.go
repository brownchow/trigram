package trigram

import "sort"

type Trigram uint32

type docList []int

func (d docList) Len() int           { return len(d) }
func (d docList) Less(i, j int) bool { return d[i] < d[j] }
func (d docList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

//IndexResult The trigram indexing result include all Documents IDs and its Frequency in that document
// 索引结果包含所有文档ID和他出现的次数
type IndexResult struct {
	// save all trigram mapping docID
	DocIDs map[int]bool

	// **文档Id**和他出现的次数
	// save all trigram appear time for trigram deletion
	Freq map[int]int
}

// ExtractStringToTrigram Extract one string to trigram list, Note the Trigram is a unit32 for ascii code
// 把输入字符串转成 Trigram 数组
func ExtractStringToTrigram(str string) []Trigram {
	if len(str) == 0 {
		return nil
	}
	var result []Trigram
	for i := 0; i < len(str)-2; i++ {
		trigram := Trigram(uint32(str[i])<<16 | uint32(str[i+1])<<8 | uint32(str[i+2]))
		result = append(result, trigram)
	}
	return result
}

//TrigramIndex 结构体，类似Java 里的class
type TrigramIndex struct {
	// to store all current trigram indexing result
	TrigramMap map[Trigram]IndexResult

	// 类似MySQL的主键ID
	// it represent and document incremental index
	maxDocID int

	// it include currently all the doc list, it will be used when query string length less than 3
	docIDsMap map[int]bool
}

// NewTrigramIndex 新建一个结构体
func NewTrigramIndex() *TrigramIndex {
	t := new(TrigramIndex)
	t.TrigramMap = make(map[Trigram]IndexResult)
	t.docIDsMap = make(map[int]bool)
	return t
}

// Add new Docoument into this trigram index
func (t *TrigramIndex) Add(doc string) int {
	// 新增一条记录
	newDocID := t.maxDocID + 1
	// 1、把输入字符串转成 trigram 数组
	trigrams := ExtractStringToTrigram(doc)
	// 2、遍历trigram数组
	for _, trigram := range trigrams {
		var mapResult IndexResult
		var exist bool
		// 3.1 如果trigram没对应的IndexResult不在的话，新建一个
		if mapResult, exist = t.TrigramMap[trigram]; !exist {
			// New doc ID handle
			mapResult = IndexResult{}
			mapResult.DocIDs = make(map[int]bool)
			mapResult.Freq = make(map[int]int)
			mapResult.DocIDs[newDocID] = true
			mapResult.Freq[newDocID] = 1
		} else {
			// 3.2 trigram对应的IndexResult已经存在，频率+1
			// trigram already exist on this doc
			if _, docExist := mapResult.DocIDs[newDocID]; docExist {
				mapResult.Freq[newDocID] = mapResult.Freq[newDocID] + 1
			} else {
				mapResult.DocIDs[newDocID] = true
				mapResult.Freq[newDocID] = 1
			}
		}
		t.TrigramMap[trigram] = mapResult
	}
	//4、返回最新的文档ID
	t.maxDocID = newDocID
	t.docIDsMap[newDocID] = true
	return newDocID
}

// Delete delete a doc for this trigram Indexing
// 删除文档和它对应的文档ID
func (t *TrigramIndex) Delete(doc string, docID int) {
	// 1、将输入字符串转成 trigram数组
	trigrams := ExtractStringToTrigram(doc)
	// 2、遍历trigram数组
	for _, trigram := range trigrams {
		//3.1、 trigram对应的IndexResult已存在
		if indexResult, exist := t.TrigramMap[trigram]; exist {
			if freq, docExist := indexResult.Freq[docID]; docExist && freq > 1 {
				indexResult.Freq[docID] = indexResult.Freq[docID] - 1
			} else {
				// need remove trigram from such docID
				delete(indexResult.Freq, docID)
				delete(indexResult.DocIDs, docID)
			}

			if len(indexResult.DocIDs) == 0 {
				// this indexResult become empty, remove this
				delete(t.TrigramMap, trigram)
				// check if some doc id has no trigram remove
			} else {
				// update back since there still other doc id exist
				t.TrigramMap[trigram] = indexResult
			}
		} else {
			// trigram not exist in map, leave
			return
		}
	}
}

// IntersectTwoMap intersect Two map
func IntersectTwoMap(IDsA, IDsB map[int]bool) map[int]bool {
	var retIDs map[int]bool   // for traversal it is smaller one
	var checkIDs map[int]bool // for checking it is bigger one
	if len(IDsA) >= len(IDsB) {
		retIDs = IDsB
		checkIDs = IDsA
	} else {
		retIDs = IDsA
		checkIDs = IDsB
	}
	for id := range retIDs {
		if _, exists := checkIDs[id]; !exists {
			delete(checkIDs, id)
		}
	}
	return retIDs
}

//Query a target string to return the doc ID
// 根据字符串查询文档
func (t *TrigramIndex) Query(doc string) docList {
	//1、把输入字符串转成trigram数组
	trigrams := ExtractStringToTrigram(doc)
	if len(trigrams) == 0 {
		return t.getAllDocIDs()
	}
	//2、拿第一个trigram去Trigram中查询 IndexResult，拿其中的 DocIDs map[int]bool
	// find first trigram as base for intersect
	retObj, exist := t.TrigramMap[trigrams[0]]
	if !exist {
		return nil
	}
	retIDs := retObj.DocIDs

	// 3、从第2个Trigram开始，以此与第1个 trigram 的结果(retIDs)做交集，将结果叠加到retIDs上
	// Remove first one and do intersect with other trigram
	trigrams = trigrams[1:]
	for _, trigram := range trigrams {
		checkObj, exist := t.TrigramMap[trigram]
		if !exist {
			return nil
		}
		checkIDs := checkObj.DocIDs
		retIDs = IntersectTwoMap(retIDs, checkIDs)
	}
	return getMapToSlice(retIDs)
}

func (t *TrigramIndex) getAllDocIDs() docList {
	return getMapToSlice(t.docIDsMap)
}

// 将DocIDs map[int]bool 转成   []int
// transfer map to slice for return result
func getMapToSlice(inMap map[int]bool) docList {
	var retSlice docList
	for k, _ := range inMap {
		retSlice = append(retSlice, k)
	}
	sort.Sort(retSlice)
	return retSlice
}
