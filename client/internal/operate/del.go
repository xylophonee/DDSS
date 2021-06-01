/**
 * @Author : gaopeihan
 * @ClassName : del.go
 * @Date : 2021/5/30 20:11
 */
package operate

import "fmt"

func (o *Operate) Del(hash string)  {
	success := o.esFile.DeleteMeta(hash)
	if success{
		fmt.Println("删除成功")
	}else {
		fmt.Println("删除错误")
	}
}
