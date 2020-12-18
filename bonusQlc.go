package bonusgo

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

//七星彩计算奖金
type BonusQlc struct {
	BaseBonus
}

/**
 * @name:计算单注奖金
 * @msg:计算单注奖金
 * @param playType 投注方式 1:普通  2 追加 lotNum 投注号码  bonusRedBall 开奖红球字符串 bonusBlueBall 篮球字符串 bonusLevel 格式化完成的奖级与开奖金额映射
 * @return: int 中奖金额 int 中奖级别
 */
func (bonus *BonusQlc) Calculation(playType int, lotNum string, bonusRedBall string, bonusBlueBall string, bonusLevel map[string][]int) (int, int) {
	var redBall []string = strings.Split(lotNum, ",")
	var bonusBall []string = strings.Split(bonusRedBall, ",")
	//数组求交集
	var redNum int = bonus.BallIntersect(bonusBall, redBall)

	var blueNum int = 0

	for _, ball := range redBall {
		if strings.EqualFold(ball, bonusBlueBall) {
			blueNum = 1
		}
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
func (bonus *BonusQlc) FormateOpenNum() ([]string, error) {
	var OpenBonusArr []string = make([]string, 0)
	if strings.Count(bonus.openNum, "-") != 1 {
		return OpenBonusArr, errors.New("开奖号码错误")
	}

	if strings.Count(bonus.openNum, ",") != 6 {
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
func (bonus *BonusQlc) SplitTicket(ticket Ticket) []string {
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
func (bonus *BonusQlc) spliteDuplicate(ticket Ticket) []string {
	var allBall []string = strings.Split(ticket.LotNum, ",")
	redZhuHe := NewZuheString(allBall, bonus.redBall.Min).ZuheResults()
	var indexRedBall int = len(redZhuHe)
	var redBallArr []string = make([]string, indexRedBall)
	for i := 0; i < indexRedBall; i++ {
		redBallArr[i] = strings.Join(redZhuHe[i], ",")
	}
	return redBallArr
}

//胆拖拆单
//01,02|03,04,05,06-01|02,03
func (bonus *BonusQlc) spliteDantuo(ticket Ticket) []string {
	var redDanTuo []string = strings.Split(ticket.LotNum, "|") //红球必有胆拖
	var redDanArr []string = strings.Split(redDanTuo[0], ",")
	var redTuoArr []string = strings.Split(redDanTuo[1], ",")
	var redDanNum int = len(redDanArr) //红球胆码数量

	var redComuNum int = bonus.redBall.Min - redDanNum //红球参与组合的数量  5 - 胆码数量

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

	return redBallArr

}

/**
 * @name:实例化大乐透的算奖结构体
 * @msg:把全部的胆拖票拆成单式，按照单式计算奖金
 * @param nil
 * @return: LotteryBonusInterface 接口(BonusDlt实现了此接口)
 */
func NewBonusQlc() LotteryBonusInterface {
	return &BonusQlc{
		BaseBonus{
			bonus: []Bonus{
				Bonus{Grade: 1, Redball: 7, BlueBall: 0, Monye: 0, Types: 1},
				Bonus{Grade: 2, Redball: 6, BlueBall: 1, Monye: 0, Types: 1},
				Bonus{Grade: 3, Redball: 6, BlueBall: 0, Monye: 0, Types: 1},
				Bonus{Grade: 4, Redball: 5, BlueBall: 1, Monye: 20000, Types: 1},
				Bonus{Grade: 5, Redball: 5, BlueBall: 0, Monye: 5000, Types: 1},
				Bonus{Grade: 6, Redball: 4, BlueBall: 1, Monye: 1000, Types: 1},
				Bonus{Grade: 7, Redball: 4, BlueBall: 0, Monye: 500, Types: 1},
			},
			redBall:   NumBall{Min: 7, Max: 30},
			blueBall:  NumBall{Min: 0, Max: 0},
			maxLervel: 7,
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
	RegisterBonusMethod("qlc", NewBonusQlc)
}
