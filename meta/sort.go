package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByUploadTime []FileMeta

//首先定义了一个时间格式的常量 baseFormat，格式为 "2006-01-02 15:04:05"。这个格式将用于解析 FileMeta 结构体中的 UploadAt 字段。
//然后定义了一个名为 ByUploadTime 的类型,它是 FileMeta 结构体切片的别名。这个类型主要用于对 FileMeta 切片进行排序。

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//这两个方法是实现了 sort.Interface 接口的必要方法:
//Len() 方法返回切片的长度。
//Swap() 方法交换切片中索引为 i 和 j 的两个元素。

func (a ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}

//使用 time.Parse() 函数将 FileMeta 结构体的 UploadAt 字段解析成 time.Time 类型的 iTime 和 jTime。
//然后比较这两个时间的纳秒级时间戳(UnixNano())。
//如果 iTime 的时间戳大于 jTime 的时间戳,则返回 true，表示 i 对应的 FileMeta 应该排在 j 对应的 FileMeta 之前。

/*总的来说,这个工具包提供了一种按照 FileMeta 结构体的 UploadAt
字段进行排序的方式。开发者可以直接使用 ByUploadTime 类型,并调用 sort.Sort() 函数对 FileMeta 切片进行排序。这可以在需要按照文
件上传时间展示列表的场景中很方便地使用。
*/
