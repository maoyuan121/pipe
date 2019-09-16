package util

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 分页信息
type Pagination struct {
	CurrentPageNum  int    `json:"currentPageNum"`  // 当前是那一页
	PageSize        int    `json:"pageSize"`        // 一页多少条数据
	PageCount       int    `json:"pageCount"`       // 总页数
	WindowSize      int    `json:"windowSize"`      // 最多显示多少个页码按钮
	RecordCount     int    `json:"recordCount"`     // 总记录数
	PageNums        []int  `json:"pageNums"`        // 页码按钮集合
	NextPageNum     int    `json:"nextPageNum"`     // 下一页的页码
	PreviousPageNum int    `json:"previousPageNum"` // 上一页的页码
	FirstPageNum    int    `json:"firstPageNum"`    // 第一页
	LastPageNum     int    `json:"lastPageNum"`     // 最后一页的页码
	PageURL         string `json:"pageURL"`         // todo
}

// 获取指定的页码
// 从 querystring 的 p 参数中获取
func GetPage(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Query("p"))
	if 1 > ret {
		ret = 1
	}

	return ret
}

// 实例化一个分页信息
func NewPagination(currentPageNum, pageSize, windowSize, recordCount int) *Pagination {
	pageCount := int(math.Ceil(float64(recordCount) / float64(pageSize)))

	previousPageNum := currentPageNum - 1
	if 1 > previousPageNum {
		previousPageNum = 0
	}
	nextPageNum := currentPageNum + 1
	if nextPageNum > pageCount {
		nextPageNum = 0
	}

	pageNums := paginate(currentPageNum, pageSize, pageCount, windowSize)
	firstPageNum := 0
	lastPageNum := 0
	if 0 < len(pageNums) {
		firstPageNum = pageNums[0]
		lastPageNum = pageNums[len(pageNums)-1]
	}

	return &Pagination{
		CurrentPageNum:  currentPageNum,
		NextPageNum:     nextPageNum,
		PreviousPageNum: previousPageNum,
		PageSize:        pageSize,
		PageCount:       pageCount,
		WindowSize:      windowSize,
		RecordCount:     recordCount,
		PageNums:        pageNums,
		FirstPageNum:    firstPageNum,
		LastPageNum:     lastPageNum,
	}
}

// 获取页码按钮集合
func paginate(currentPageNum, pageSize, pageCount, windowSize int) []int {
	var ret []int

	if pageCount < windowSize {
		for i := 0; i < pageCount; i++ {
			ret = append(ret, i+1)
		}
	} else {
		first := currentPageNum + 1 - windowSize/2
		if first < 1 {
			first = 1
		}
		if first+windowSize > pageCount {
			first = pageCount - windowSize + 1
		}
		for i := 0; i < windowSize; i++ {
			ret = append(ret, first+i)
		}
	}

	return ret
}
