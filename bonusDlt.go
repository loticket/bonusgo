package bonusgo

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type BonusDlt struct {
	BaseBonus
}

/**
 * @name:计算一张票的奖金
 * @msg:计算单注奖金，可以设定完成开奖号码，和开奖级别处理以后，通过设置不同的ticket计算同期号的所有奖金
 * @param bonusOpen 开奖号码数组 []string{红球字符串，篮球字符串}(可调用 FormateOpenNum()函数获取)
 * @param bonusLevel 开奖信息  map[string][]int map[红球个数_篮球个数_开奖方式][]int{奖级，奖金}(可调用 BonusLevelFormate()函数)
 * @return: int64 税前金额  int64 税后金额  int 是否是大额奖金  []int 中奖注数 string 票的唯一标示
 */
func (bonus *BonusDlt) CalculatePrize(ticket Ticket, bonusOpen []string, bonusLevel map[string][]int) (int64, int64, int, []int, string) {
	//格式化中奖
	var taxBeforeMoneyAll int64 = 0                     //税前金额
	var taxAfterMoneyAll int64 = 0                      //税后金额
	var levels []int = []int{0, 0, 0, 0, 0, 0, 0, 0, 0} //中奖注数
	var big int = 0
	//获取拆票信息
	var ticketAll []string = bonus.SplitTicket(ticket)

	for i := 0; i < len(ticketAll); i++ {
		grad, money := bonus.calculation(ticket.PlayType, ticketAll[i], bonusOpen[0], bonusOpen[1], bonusLevel)
		if grad == 0 {
			continue
		}
		taxBeforeMoneyAll += int64(money * ticket.Multiple)
		big = 2
		if grad < 3 && money > 1000000 {
			taxAfterMoneyAll += int64((money * 8 / 10) * ticket.Multiple)
			big = 1
		} else {
			taxAfterMoneyAll += int64(money * ticket.Multiple)
		}

		levels[grad-1]++
	}

	return taxBeforeMoneyAll, taxAfterMoneyAll, big, levels, ticket.Tid
}

/**
 * @name:计算单注奖金
 * @msg:计算单注奖金
 * @param playType 投注方式 1:普通  2 追加 lotNum 投注号码  bonusRedBall 开奖红球字符串 bonusBlueBall 篮球字符串 bonusLevel 格式化完成的奖级与开奖金额映射
 * @return: int 中奖金额 int 中奖级别
 */
func (bonus *BonusDlt) calculation(playType int, lotNum string, bonusRedBall string, bonusBlueBall string, bonusLevel map[string][]int) (int, int) {
	var allBall []string = strings.Split(lotNum, "-")
	var redBall []string = strings.Split(allBall[0], ",")
	var blueBall []string = strings.Split(allBall[1], ",")
	var redNum int = bonus.BallNumInBonus(redBall, bonusRedBall, bonus.redBall.Min)
	var blueNum int = bonus.BallNumInBonus(blueBall, bonusBlueBall, bonus.blueBall.Min)

	var money int = 0 //中奖金额

	if redNum == 0 && blueNum == 0 {
		return 0, 0
	}

	//直接计算出来奖金
	var keys string = strconv.Itoa(redNum) + "_" + strconv.Itoa(blueNum)

	var normal string = keys + "_1"

	levelCal, ok := bonusLevel[normal]

	if !ok {
		return 0, 0
	}

	money += levelCal[1]

	//如果是追加这加上普通奖金
	if playType == 1 || (playType == 2 && levelCal[0] != 1 && levelCal[0] != 2) {
		return levelCal[0], levelCal[1]
	}

	var addBonus string = keys + "_2"

	levelCaAdd, oks := bonusLevel[addBonus]
	if !oks {
		return 0, 0
	}

	money += levelCaAdd[1]

	return levelCal[0], money
}

/**
 * @name:简单验证开奖号码并且格式化开奖号码
 * @msg:在计算开奖之前一定要验证开奖信息
 * @param nil
 * @return: []string{红球字符串,篮球字符串} error 错误提示
 */
func (bonus *BonusDlt) FormateOpenNum() ([]string, error) {
	var OpenBonusArr []string = make([]string, 0)
	if strings.Count(bonus.openNum, "-") != 1 {
		return OpenBonusArr, errors.New("开奖号码错误")
	}

	if strings.Count(bonus.openNum, ",") != 5 {
		return OpenBonusArr, errors.New("开奖号码错误")
	}

	return strings.Split(bonus.openNum, "-"), nil
}

/**
 * @name:单式，复式，胆拖拆票
 * @msg:把全部的票拆成单式，按照单式计算奖金
 * @param nil
 * @return: []string{单式号码字符串,单式号码字符串}
 */
func (bonus *BonusDlt) SplitTicket(ticket Ticket) []string {
	switch ticket.LotType {
	case 1:
		return []string{ticket.LotNum}
	case 2:
		return bonus.spliteDuplicate(ticket)
	case 3:
		return bonus.spliteDantuo(ticket)
	default:
		return []string{}
	}
	return []string{}
}

/**
 * @name:复式拆票
 * @msg:把全部的复式票拆成单式，按照单式计算奖金
 * @param nil
 * @return: []string{单式号码字符串,单式号码字符串}
 */
func (bonus *BonusDlt) spliteDuplicate(ticket Ticket) []string {
	var allBall []string = strings.Split(ticket.LotNum, "-")
	var redBall []string = strings.Split(allBall[0], ",")
	var blueBall []string = strings.Split(allBall[1], ",")
	sort.Strings(redBall)
	sort.Strings(blueBall)
	var redBallArr []string
	var blueBallArr []string
	var redBallNum int = len(redBall)   //红球数量
	var blueBallNum int = len(blueBall) //篮球数量
	var indexRedBall int = 0
	var indexBlueBall int = 0
	if redBallNum == 5 {
		indexRedBall = 1
		redBallArr = make([]string, 1)
		redBallArr[0] = allBall[0]
	} else {
		redZhuHe := NewZuheString(redBall, 5).ZuheResults()
		indexRedBall = len(redZhuHe)
		redBallArr = make([]string, indexRedBall)
		for i := 0; i < indexRedBall; i++ {
			redBallArr[i] = strings.Join(redZhuHe[i], ",")
		}
	}

	if blueBallNum == 2 {
		indexBlueBall = 1
		blueBallArr = make([]string, 1)
		blueBallArr[0] = allBall[1]
	} else {
		blueZhuHe := NewZuheString(blueBall, 2).ZuheResults()
		indexBlueBall = len(blueZhuHe)
		blueBallArr = make([]string, indexBlueBall)
		for i := 0; i < indexBlueBall; i++ {
			blueBallArr[i] = strings.Join(blueZhuHe[i], ",")
		}
	}

	//组合最后的结果
	var resZuHe []string = make([]string, 0)

	for j := 0; j < indexRedBall; j++ {
		for k := 0; k < indexBlueBall; k++ {
			resZuHe = append(resZuHe, (redBallArr[j] + "-" + blueBallArr[k]))
		}
	}

	return resZuHe

}

/**
 * @name:胆拖拆票
 * @msg:把全部的胆拖票拆成单式，按照单式计算奖金
 * @param nil
 * @return: []string{单式号码字符串,单式号码字符串}
 */
func (bonus *BonusDlt) spliteDantuo(ticket Ticket) []string {
	var allBall []string = strings.Split(ticket.LotNum, "-")
	var redDanTuo []string = strings.Split(allBall[0], "|") //红球必有胆拖
	var redDanArr []string = strings.Split(redDanTuo[0], ",")
	var redTuoArr []string = strings.Split(redDanTuo[1], ",")
	var redDanNum int = len(redDanArr) //红球胆码数量
	var blueDanArr []string
	var blueTuoArr []string
	var blueDanNum int = 0
	var blueTuoNum int = 0

	if strings.Contains(allBall[1], "|") { //如果有胆拖
		var blueDanTuo []string = strings.Split(allBall[1], "|") //篮球胆拖
		blueDanArr = strings.Split(blueDanTuo[0], ",")
		blueDanNum = len(blueDanArr)
		blueTuoArr = strings.Split(blueDanTuo[1], ",")
		blueTuoNum = len(blueTuoArr)
	} else {
		blueDanArr = make([]string, 0)
		blueDanNum = 0
		blueTuoArr = strings.Split(allBall[1], ",") //篮球胆拖
		blueTuoNum = len(blueTuoArr)
	}

	var redComuNum int = 5 - redDanNum //红球参与组合的数量  5 - 胆码数量

	var blueComuNum int = 2 - blueDanNum //篮球参与组合的数量  2——胆码数
	sort.Strings(redTuoArr)
	sort.Strings(blueTuoArr)

	//组合单式
	var redBallArr []string = make([]string, 0)
	var blueBallArr []string = make([]string, 0)
	//组合红球
	redZhuHe := NewZuheString(redTuoArr, redComuNum).ZuheResults()
	var indexRedBall int = len(redZhuHe)
	redBallArr = make([]string, indexRedBall)
	for i := 0; i < indexRedBall; i++ {
		redBallArr[i] = strings.Join(redDanArr, ",") + "," + strings.Join(redZhuHe[i], ",")
	}

	var indexBlueBall int = 1
	//组合篮球
	if blueDanNum == 0 && blueTuoNum == 2 {
		blueBallArr = append(blueBallArr, allBall[1])
	} else if blueDanNum == 0 && blueTuoNum > 2 {
		blueZhuHe := NewZuheString(blueTuoArr, blueComuNum).ZuheResults()
		indexBlueBall = len(blueZhuHe)
		for i := 0; i < indexBlueBall; i++ {
			blueBallArr = append(blueBallArr, strings.Join(blueZhuHe[i], ","))
		}
	} else if blueDanNum == 1 {
		blueZhuHe := NewZuheString(blueTuoArr, blueComuNum).ZuheResults()
		indexBlueBall = len(blueZhuHe)
		for i := 0; i < indexBlueBall; i++ {
			blueBallArr = append(blueBallArr, strings.Join(blueDanArr, ",")+","+strings.Join(blueZhuHe[i], ","))
		}
	}

	//组合最后的结果
	var resZuHe []string = make([]string, 0)

	for j := 0; j < indexRedBall; j++ {
		for k := 0; k < indexBlueBall; k++ {
			resZuHe = append(resZuHe, (redBallArr[j] + "-" + blueBallArr[k]))
		}
	}

	return resZuHe
}

/**
 * @name:实例化大乐透的算奖结构体
 * @msg:把全部的胆拖票拆成单式，按照单式计算奖金
 * @param nil
 * @return: LotteryBonusInterface 接口(BonusDlt实现了此接口)
 */
func NewBonusDlt() LotteryBonusInterface {
	return &BonusDlt{
		BaseBonus{
			bonus: []Bonus{
				Bonus{Grade: 1, Redball: 5, BlueBall: 2, Monye: 0, Types: 1}, //types 1 普通  2 追加
				Bonus{Grade: 1, Redball: 5, BlueBall: 2, Monye: 0, Types: 2}, //types 1 普通  2 追加
				Bonus{Grade: 2, Redball: 5, BlueBall: 1, Monye: 0, Types: 1},
				Bonus{Grade: 2, Redball: 5, BlueBall: 1, Monye: 0, Types: 2},
				Bonus{Grade: 3, Redball: 5, BlueBall: 0, Monye: 1000000, Types: 1},
				Bonus{Grade: 4, Redball: 4, BlueBall: 2, Monye: 300000, Types: 1},
				Bonus{Grade: 5, Redball: 4, BlueBall: 1, Monye: 30000, Types: 1},
				Bonus{Grade: 6, Redball: 3, BlueBall: 2, Monye: 20000, Types: 1},
				Bonus{Grade: 7, Redball: 4, BlueBall: 0, Monye: 10000, Types: 1},
				Bonus{Grade: 8, Redball: 3, BlueBall: 1, Monye: 1500, Types: 1},
				Bonus{Grade: 8, Redball: 2, BlueBall: 2, Monye: 1500, Types: 1},
				Bonus{Grade: 9, Redball: 3, BlueBall: 0, Monye: 500, Types: 1},
				Bonus{Grade: 9, Redball: 1, BlueBall: 2, Monye: 500, Types: 1},
				Bonus{Grade: 9, Redball: 2, BlueBall: 1, Monye: 500, Types: 1},
				Bonus{Grade: 9, Redball: 0, BlueBall: 2, Monye: 500, Types: 1},
			},
			redBall:  NumBall{Min: 5, Max: 35},
			blueBall: NumBall{Min: 2, Max: 12},
		},
	}
}

/**
 * @name:实例化大乐透算奖结构体
 * @msg:实例化大乐透算奖结构体
 * @param nil
 * @return: func() LotteryBonusInterface 接口(BonusDlt实现了此接口)
 */
func init() {
	RegisterBonusMethod("dlt", NewBonusDlt)
}
