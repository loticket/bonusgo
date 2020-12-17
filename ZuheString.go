package bonusgo

type ZuheString struct {
	Nums []string //排列的数字
	Nnum int      //总个数
	Mnum int      //取出的数量

}

func NewZuheString(nums []string, mnum int) *ZuheString {
	return &ZuheString{
		Nums: nums,
		Mnum: mnum,
	}
}

func (this *ZuheString) ZuheResults() [][]string {
	this.Nnum = len(this.Nums)
	return this.findNumsByIndexs()
}

//组合算法(从nums中取出m个数)
func (this *ZuheString) zuheResult() [][]int {
	if this.Mnum < 1 || this.Mnum > this.Nnum {
		return [][]int{}
	}

	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, this.mathZuhe())
	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexs := make([]int, this.Nnum)
	for i := 0; i < this.Nnum; i++ {
		if i < this.Mnum {
			indexs[i] = 1
		} else {
			indexs[i] = 0
		}
	}

	//第一个结果
	result = this.addTo(result, indexs)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < this.Nnum-1; i++ {
			if indexs[i] == 1 && indexs[i+1] == 0 {
				find = true

				indexs[i], indexs[i+1] = 0, 1
				if i > 1 {
					this.moveOneToLeft(indexs[:i])
				}
				result = this.addTo(result, indexs)

				break
			}
		}

		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}

	return result
}

//数学方法计算组合数(从n中取m个数)
func (this *ZuheString) mathZuhe() int {
	return Combination(this.Nnum, this.Mnum)
}

//将ele复制后添加到arr中，返回新的数组
func (this *ZuheString) addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)

	return arr
}

func (this *ZuheString) moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}

	//将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

//根据索引号数组得到元素数组
func (this *ZuheString) findNumsByIndexs() [][]string {
	var indexs [][]int = this.zuheResult()
	if len(indexs) == 0 {
		return [][]string{}
	}

	result := make([][]string, len(indexs))

	for i, v := range indexs {
		line := make([]string, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, this.Nums[j])
			}
		}
		result[i] = line
	}

	return result
}

//组合计算公式
func Combination(n int, r int) int {
	return int(Factorial(n) / (Factorial(n-r) * Factorial(r)))
}

//阶乘计算
func Factorial(n int) int64 {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	var result int64 = 1
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}
