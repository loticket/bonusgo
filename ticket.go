package bonusgo

//数字彩投注内容
type Ticket struct {
	Tid      string `json:"tid" binding:"required"`
	PlayType int    `json:"playtype" binding:"required"` //子玩法
	LotType  int    `json:"lottype" binding:"required"`  //选号方式
	LotNum   string `json:"lotnum" binding:"required"`   //投注号码
	Money    int    `json:"money" binding:"required"`    //投注金额
	BetNum   int    `json:"betnum" binding:"required"`   //注数
	Multiple int    `json:"multiple" binding:"required"` //倍数
}

//开奖数据
type OpenBonus struct {
	LotId    string  `json:"lotid" binding:"lotid"`       //彩种编码
	OpenNum  string  `json:"opennum" binding:"opennum"`   //开奖号码
	OpenInfo []Bonus `json:"openinfo" binding:"openinfo"` //开奖号码
}

//每一个彩种计算奖金的接口，都必须要实现此接口
type LotteryBonusInterface interface {
	SetTicket(Ticket)                    //票信息
	SetOpenBonus(string)                 //string 开奖号码
	SetBonusMoney([]int)                 //[]int 为每个等级的奖金
	GetBonusMoney() []Bonus              //获取奖金数组
	FormateOpenNum() ([]string, error)   //格式化开奖号码
	BonusLevelFormate() map[string][]int //格式化开奖级别信息
	//CalculatePrize(Ticket, []string, map[string][]int) (int64, int64, int, []int, string) //计算奖金并返回
	Calculation(int, string, string, string, map[string][]int) (int, int)
	SplitTicket(Ticket) []string //拆票
	GetmaxLervel() int
}

//中奖级别定义
type Bonus struct {
	Grade    int //中奖等级
	Redball  int //红球个数
	BlueBall int //篮球个数
	Monye    int //中奖金额
	Types    int //奖金类型  1 普通   2 其他
}

//定义球的数量
type NumBall struct {
	Min int //最小数
	Max int //最大数
}

//结果定义成结构体
type TicketCalResult struct {
	Tid               string
	TaxBeforeMoneyAll int64
	TaxAfterMoneyAll  int64
	Big               int
	Levels            []int
}
