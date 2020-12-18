package bonusgo

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type BonusSsq struct {
	BaseBonus
}

/**
 * @name:计算单注奖金
 * @msg:计算单注奖金
 * @param playType 投注方式 1:普通  2 追加 lotNum 投注号码  bonusRedBall 开奖红球字符串 bonusBlueBall 篮球字符串 bonusLevel 格式化完成的奖级与开奖金额映射
 * @return: int 中奖金额 int 中奖级别
 */
func (bonus *BonusSsq) Calculation(playType int, lotNum string, bonusRedBall string, bonusBlueBall string, bonusLevel map[string][]int) (int, int) {
	var allBall []string = strings.Split(lotNum, "-")
	var redBall []string = strings.Split(allBall[0], ",")
	var redNum int = bonus.BallNumInBonus(redBall, bonusRedBall, bonus.redBall.Min)
	var blueNum int = 0

	if strings.EqualFold(allBall[1], bonusBlueBall) {
		blueNum = 1
	}

	if redNum == 0 && blueNum == 0 {
		return 0, 0
	}

	//直接计算出来奖金
	var keys string = strconv.Itoa(redNum) + "_" + strconv.Itoa(blueNum) + "_" + strconv.Itoa(playType)

	levelCal, ok := bonusLevel[keys]

	if !ok {
		return 0, 0
	}

	return levelCal[0], levelCal[1]
}

/**
 * @name:简单验证开奖号码并且格式化开奖号码
 * @msg:在计算开奖之前一定要验证开奖信息
 * @param nil
 * @return: []string{红球字符串,篮球字符串} error 错误提示
 */
func (bonus *BonusSsq) FormateOpenNum() ([]string, error) {
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
func (bonus *BonusSsq) SplitTicket(ticket Ticket) []string {
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
func (bonus *BonusSsq) spliteDuplicate(ticket Ticket) []string {
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
	if redBallNum == 6 {
		indexRedBall = 1
		redBallArr = make([]string, 1)
		redBallArr[0] = allBall[0]
	} else {
		redZhuHe := NewZuheString(redBall, bonus.redBall.Min).ZuheResults()
		indexRedBall = len(redZhuHe)
		redBallArr = make([]string, indexRedBall)
		for i := 0; i < indexRedBall; i++ {
			redBallArr[i] = strings.Join(redZhuHe[i], ",")
		}
	}

	if blueBallNum == 1 {
		indexBlueBall = 1
		blueBallArr = make([]string, 1)
		blueBallArr[0] = allBall[1]
	} else {
		blueZhuHe := NewZuheString(blueBall, 1).ZuheResults()
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
func (bonus *BonusSsq) spliteDantuo(ticket Ticket) []string {
	var allBall []string = strings.Split(ticket.LotNum, "-")
	var redDanTuo []string = strings.Split(allBall[0], "|") //红球必有胆拖
	var redDanArr []string = strings.Split(redDanTuo[0], ",")
	var redTuoArr []string = strings.Split(redDanTuo[1], ",")
	var redDanNum int = len(redDanArr)                        //红球胆码数量
	var blueBallArr []string = strings.Split(allBall[1], ",") //篮球
	var blueBallNum int = len(blueBallArr)

	var redComuNum int = bonus.redBall.Min - redDanNum //红球参与组合的数量  6 - 胆码数量

	sort.Strings(redTuoArr)

	//组合单式
	var redBallArr []string = make([]string, 0)

	//组合红球
	redZhuHe := NewZuheString(redTuoArr, redComuNum).ZuheResults()
	var indexRedBall int = len(redZhuHe)
	redBallArr = make([]string, indexRedBall)
	for i := 0; i < indexRedBall; i++ {
		redBallArr[i] = strings.Join(redDanArr, ",") + "," + strings.Join(redZhuHe[i], ",")
	}

	//组合最后的结果
	var resZuHe []string = make([]string, 0)

	for j := 0; j < indexRedBall; j++ {
		for k := 0; k < blueBallNum; k++ {
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
func NewBonusSsq() LotteryBonusInterface {
	return &BonusSsq{
		BaseBonus{
			bonus: []Bonus{
				Bonus{Grade: 1, Redball: 6, BlueBall: 1, Monye: 0, Types: 1},
				Bonus{Grade: 2, Redball: 6, BlueBall: 0, Monye: 0, Types: 1},
				Bonus{Grade: 3, Redball: 5, BlueBall: 1, Monye: 300000, Types: 1},
				Bonus{Grade: 4, Redball: 5, BlueBall: 0, Monye: 20000, Types: 1},
				Bonus{Grade: 4, Redball: 4, BlueBall: 1, Monye: 20000, Types: 1},
				Bonus{Grade: 5, Redball: 4, BlueBall: 0, Monye: 1000, Types: 1},
				Bonus{Grade: 5, Redball: 3, BlueBall: 1, Monye: 1000, Types: 1},
				Bonus{Grade: 6, Redball: 2, BlueBall: 1, Monye: 500, Types: 1},
				Bonus{Grade: 6, Redball: 1, BlueBall: 1, Monye: 500, Types: 1},
				Bonus{Grade: 6, Redball: 0, BlueBall: 1, Monye: 500, Types: 1},
			},
			redBall:   NumBall{Min: 6, Max: 35},
			blueBall:  NumBall{Min: 1, Max: 12},
			maxLervel: 6,
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
	RegisterBonusMethod("ssq", NewBonusSsq)
}
