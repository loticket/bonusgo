package bonusgo

import (
	"strconv"
	"strings"
)

type BaseBonus struct {
	openNum   string  //开奖号码
	ticket    Ticket  //票详情
	bonus     []Bonus //奖级别数组
	redBall   NumBall //红球个数
	blueBall  NumBall //篮球个数
	maxLervel int
}

//传票进入
func (bonus *BaseBonus) SetTicket(ticket Ticket) {
	bonus.ticket = ticket
}

//传入开奖信息
func (bonus *BaseBonus) SetOpenBonus(openNum string) {
	bonus.openNum = openNum
}

/**
 * @name:设定开奖级别的奖金
 * @msg:Bonus 包含几等奖，需要命中球的个数 中奖金额
 * @param []int 整数数组，金额以分为单位，[]int数组的金额需要与[]bonus对应
 * @return: void
 */
func (bonus *BaseBonus) SetBonusMoney(money []int) {
	for i := 0; i < len(money); i++ {
		bonus.bonus[i].Monye = money[i]
	}
}

/**
 * @name:获取开奖级别的相信
 * @msg:Bonus 包含几等奖，需要命中球的个数 中奖金额
 * @param nil
 * @return: []Bonus
 */
func (bonus *BaseBonus) GetBonusMoney() []Bonus {
	return bonus.bonus
}

func (bonus *BaseBonus) GetmaxLervel() int {
	return bonus.maxLervel
}

/**
 * @name:格式化开奖信息
 * @msg:格式化开奖信息，把红球的个数+篮球个数(若没有蓝球这填0)+开奖方式
 * @param nil
 * @return: []int{级别,奖金}
 */
func (bonus *BaseBonus) BonusLevelFormate() map[string][]int {
	var bonusLerevel map[string][]int = make(map[string][]int, 0)
	for _, v := range bonus.bonus {
		var mapkey string = strconv.Itoa(v.Redball) + "_" + strconv.Itoa(v.BlueBall) + "_" + strconv.Itoa(v.Types)
		bonusLerevel[mapkey] = []int{v.Grade, v.Monye}
	}
	return bonusLerevel
}

/**
 * @name:中奖球包含投注球的个数
 * @msg:中奖球包含投注球的个数，判断投注号码和开奖号码相同的个数
 * @param []string 投注号码数组，bonusBlueBall 开奖号码字符串 ballNum 查找次数
 * @return: int 包含的个数
 */
func (bonus *BaseBonus) BallNumInBonus(blueBall []string, bonusBlueBall string, ballNum int) int {
	var i int = 0
	for j := 0; j < ballNum; j++ {
		if strings.Contains(bonusBlueBall, blueBall[j]) {
			i++
		}
	}
	return i
}

/**
 * @name:两个字符串数组求交集
 * @msg:两个字符串数组求交集，判断投注号码和开奖号码相同的个数
 * @param []string 字符串数组
 * @return: int 交集的个数
 */
func (bonus *BaseBonus) BallIntersect(nums1 []string, nums2 []string) int {
	map1 := map[string]int{}
	for _, v := range nums1 {
		map1[v] += 1
	}
	var i int = 0
	for _, v := range nums2 {
		if map1[v] > 0 {
			i++
			map1[v] -= 1
		}
	}
	return i
}
