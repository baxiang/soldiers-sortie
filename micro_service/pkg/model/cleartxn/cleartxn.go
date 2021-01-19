package cleartxn

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//Mysql Select
const (
	SelectInsTxnResp = "KEY_RSP, MCHNT_CD, CARD_ACCPTR_NM, TRANS_DT, MA_SETTLE_DT, TRANS_MT, " +
		"MA_TRANS_CD, FWD_INS_ID_CD, TRANS_AT, PRI_ACCT_NO, ISS_INS_ID_CD,  TERM_ID, " +
		"PROD_CD , CARD_CLASS, TRANS_ST, RESP_CD"
)

//Table
const (
	clearTxnName  = "TBL_CLEAR_TXN"
	TfrTrnLogName = "TBL_TFR_TRN_LOG1"
)

type CommonModel struct {
	ID        uint      `gorm:"column:id;primary_key;auto_increment;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}

type ClearTxn struct {
	CompanyCd        string    `gorm:"column:COMPANY_CD`
	InsIDCd          string    `gorm:"column:INS_ID_CD`
	AcqInsIDCd       string    `gorm:"column:ACQ_INS_ID_CD"`
	FwdInsIDCd       string    `gorm:"column:FWD_INS_ID_CD"`
	MchtCd           string    `gorm:"column:MCHT_CD"`
	MchtName         string    `gorm:"column:MCHT_NAME"`
	MchtShortName    string    `gorm:"column:MCHT_SHORT_NAME"`
	MccCd            string    `gorm:"column:MCC_CD"`
	MccCd42          string    `gorm:"column:MCC_CD_42"`
	MccDesc          string    `gorm:"column:MCC_DESC"`
	TransDateTime    string    `gorm:"column:TRANS_DATE_TIME"`
	StlmDate         string    `gorm:"column:STLM_DATE"`
	TransKind        string    `gorm:"column:TRANS_KIND"`
	TxnDesc          string    `gorm:"column:TXN_DESC"`
	TransState       string    `gorm:"column:TRANS_STATE`
	StlmFlg          string    `gorm:"column:STLM_FLG`
	TransAmt         string    `gorm:"column:TRANS_AMT"`
	Creditcardlimit  string    `gorm:"column:CREDITCARDLIMIT"`
	CupSsn           string    `gorm:"column:CUP_SSN"`
	AuthIDResp       string    `gorm:"column:AUTHR_ID_RESP"`
	Pan              string    `gorm:"column:PAN"`
	CardKindDis      string    `gorm:"column:CARD_KIND_DIS"`
	BankCode         string    `gorm:"column:BANK_CODE"`
	BankName         string    `gorm:"column:BANK_NAME"`
	BranchCd         string    `gorm:"column:BRANCH_CD"`
	BranchNm         string    `gorm:"column:BRANCH_NM"`
	TermID           string    `gorm:"column:TERM_ID"`
	OrgTransDateTime string    `gorm:"column:ORG_TRANS_DATE_TIME"`
	OrgCupSsn        string    `gorm:"column:ORG_CUP_SSN"`
	PosEntryMode     string    `gorm:"column:POS_ENTRY_MODE"`
	RspCode          string    `gorm:"column:RSP_CODE"`
	TrueFeeMod       string    `gorm:"column:TRUE_FEE_MOD"`
	TrueFeeBi        string    `gorm:"column:TRUE_FEE_BI"`
	TrueFeeFd        string    `gorm:"column:TRUE_FEE_FD"`
	TrueFeeFfd       string    `gorm:"column:TRUE_FEE_FFD"`
	Var1             string    `gorm:"column:VAR_1"`
	Var2             string    `gorm:"column:VAR_2"`
	Var3             string    `gorm:"column:VAR_3"`
	Var4             string    `gorm:"column:VAR_4"`
	VirFeeMod        string    `gorm:"column:VIR_FEE_MOD"`
	VirFeeBi         string    `gorm:"column:VIR_FEE_BI"`
	VirFeeBd         string    `gorm:"column:VIR_FEE_BD"`
	VirFeeFd         string    `gorm:"column:VIR_FEE_FD"`
	MchtFee          string    `gorm:"column:MCHT_FEE"`
	Var5             string    `gorm:"column:VAR_5"`
	MchtVirFee       string    `gorm:"column:MCHT_VIR_FEE"`
	StandBankFee     string    `gorm:"column:STAND_BANK_FEE"`
	BankFee          string    `gorm:"column:BANK_FEE"`
	HzjgFee          string    `gorm:"column:HZJG_FEE"`
	Jgsy             string    `gorm:"column:JGSY"`
	AipFee           string    `gorm:"column:AIP_FEE"`
	MchtSetAmt       string    `gorm:"column:MCHT_SET_AMT"`
	Hzjgyfppfwf      string    `gorm:"column:HZJGYFPPFWF"`
	Jgyfppfwf        string    `gorm:"column:JGYFPPFWF"`
	Aipyfppfwf       string    `gorm:"column:AIPYFPPFWF"`
	ErrFeeIn         string    `gorm:"column:ERR_FEE_IN"`
	ErrFeeOut        string    `gorm:"column:ERR_FEE_OUT"`
	ErrCode          string    `gorm:"column:ERR_CODE"`
	JtMchtCd         string    `gorm:"column:JT_MCHT_CD"`
	ExpandOrgCd      string    `gorm:"column:EXPAND_ORG_CD"`
	SpeServInst      string    `gorm:"column:SPE_SERV_INST"`
	PropIns          string    `gorm:"column:PROP_INS"`
	ExpandOrgFee     string    `gorm:"column:EXPAND_ORG_FEE"`
	SpeServFee       string    `gorm:"column:SPE_SERV_FEE"`
	PropInsFee       string    `gorm:"column:PROP_INS_FEE"`
	ExpandOrgPp      string    `gorm:"column:EXPAND_ORG_PP"`
	SpeServPp        string    `gorm:"column:SPE_SERV_PP"`
	PropInsPp        string    `gorm:"column:PROP_INS_PP"`
	ExpandFeeIn      string    `gorm:"column:EXPAND_FEE_IN"`
	ExpandFeeOut     string    `gorm:"column:EXPAND_FEE_OUT"`
	CupIfinsideSign  string    `gorm:"column:CUP_IFINSIDE_SIGN"`
	SpChargType      string    `gorm:"column:SP_CHARG_TYPE"`
	SpChargLev       string    `gorm:"column:SP_CHARG_LEV"`
	TermSsn          string    `gorm:"column:TERM_SSN"`
	SnSsn            string    `gorm:"column:SN_SSN"`
	UpChlID          string    `gorm:"column:UP_CHL_ID"`
	ConvMchtCd       string    `gorm:"column:CONV_MCHT_CD"`
	ConvTermID       string    `gorm:"column:CONV_TERM_ID"`
	ChlTrueFee       string    `gorm:"column:CHL_TRUE_FEE"`
	ChlStdFee        string    `gorm:"column:CHL_STD_FEE"`
	ChlFeePreFlg     string    `gorm:"column:CHL_FEE_PRE_FLG"`
	SysSer           string    `gorm:"column:SYS_SER"`
	Var6             string    `gorm:"column:VAR_6"`
	QudaoFee         string    `gorm:"column:QUDAO_FEE"`
	QudaoFeeMin      string    `gorm:"column:QUDAO_FEE_MIN"`
	QudaoFeeMix      string    `gorm:"column:QUDAO_FEE_MIX"`
	QudaoFeeFd       string    `gorm:"column:QUDAO_FEE_FD"`
	InsFee           string    `gorm:"column:INS_FEE"`
	InsMyFee         string    `gorm:"column:INS_MY_FEE"`
	InsCostFee       string    `gorm:"column:INS_COST_FEE"`
	InsMyFeeAmt      string    `gorm:"column:INS_MY_FEE_AMT"`
	InsSplitFee      string    `gorm:"column:INS_SPLIT_FEE"`
	InsResFee        string    `gorm:"column:INS_RES_FEE"`
	PinpFee          string    `gorm:"column:PINP_FEE"`
	PinpFeeInf       string    `gorm:"column:PINP_FEE_INF"`
	PinpFeeTop       string    `gorm:"column:PINP_FEE_TOP"`
	PinpStat         string    `gorm:"column:PINP_STAT"`
	T0Stat           string    `gorm:"column:T0_STAT"`
	KeyRsp           string    `gorm:"column:KEY_RSP; primary_key"`
	Remark           *string   `gorm:"column:REMARK"`
	Remark1          *string   `gorm:"column:REMARK1"`
	Remark2          *string   `gorm:"column:REMARK2"`
	Remark3          *string   `gorm:"column:REMARK3"`
	Remark4          *string   `gorm:"column:REMARK4"`
	Remark5          *string   `gorm:"column:REMARK5"`
	MaSettleDt       time.Time `gorm:"column:MA_SETTLE_DT"`
	RetriRefNo       string    `gorm:"column:RETRI_REF_NO"`
	IndustryAddnInf  string    `gorm:"column:INDUSTRY_ADDN_INF"`
	SrcQid           float64   `gorm:"column:SRC_QID"`
	ProdCd           string    `gorm:"column:PROD_CD"`
}

func (c ClearTxn) TableName() string {
	return clearTxnName
}

type TfrTrnLog struct {
	TransDt              string    `gorm:"column:TRANS_DT" downI:"4"`
	TransMt              string    `gorm:"column:TRANS_MT" downI:"6"`
	SrcQid               int       `gorm:"column:SRC_QID"`
	DesQid               int       `gorm:"column:DES_QID"`
	MaTransCd            string    `gorm:"column:MA_TRANS_CD" downI:"7"`
	MaTransNm            string    `gorm:"column:MA_TRANS_NM"`
	KeyRsp               string    `gorm:"column:KEY_RSP" downI:"1"`
	KeyRevsal            string    `gorm:"column:KEY_REVSAL"`
	KeyCancel            string    `gorm:"column:KEY_CANCEL"`
	RespCd               string    `gorm:"column:RESP_CD" downI:"16"`
	TransSt              string    `gorm:"column:TRANS_ST" downI:"15"`
	MaTransSeq           int       `gorm:"column:MA_TRANS_SEQ"`
	OrigMaTransSeq       int       `gorm:"column:ORIG_MA_TRANS_SEQ"`
	OrigTransSeq         string    `gorm:"column:ORIG_TRANS_SEQ"`
	OrigTermSeq          string    `gorm:"column:ORIG_TERM_SEQ"`
	OrigTransDt          string    `gorm:"column:ORIG_TRANS_DT"`
	MaSettleDt           string    `gorm:"column:MA_SETTLE_DT" downI:"5"`
	AccessMd             string    `gorm:"column:ACCESS_MD"`
	MsgTp                string    `gorm:"column:MSG_TP"`
	PriAcctNo            string    `gorm:"column:PRI_ACCT_NO" downI"10"`
	AcctTp               string    `gorm:"column:ACCT_TP"`
	TransProcCd          string    `gorm:"column:TRANS_PROC_CD"`
	TransAt              string    `gorm:"column:TRANS_AT" downI:"9"`
	TransTdTm            string    `gorm:"column:TRANS_TD_TM"`
	TermSeq              string    `gorm:"column:TERM_SEQ"`
	AcptTransTm          string    `gorm:"column:ACPT_TRANS_TM"`
	AcptTransDt          string    `gorm:"column:ACPT_TRANS_DT"`
	MchntTp              string    `gorm:"column:MCHNT_TP"`
	PosEntryMdCd         string    `gorm:"column:POS_ENTRY_MD_CD"`
	PosCondCd            string    `gorm:"column:POS_COND_CD"`
	AcptInsIdCd          string    `gorm:"column:ACPT_INS_ID_CD"`
	FwdInsIdCd           string    `gorm:"column:FWD_INS_ID_CD" downI:"8"`
	TermId               string    `gorm:"column:TERM_ID" downI:"12"`
	MchntCd              string    `gorm:"column:MCHNT_CD" downI:"2"`
	CardAccptrNm         string    `gorm:"column:CARD_ACCPTR_NM" downI"3"`
	RetriRefNo           string    `gorm:"column:RETRI_REF_NO"`
	ReqAuthId            string    `gorm:"column:REQ_AUTH_ID"`
	TransSubcata         string    `gorm:"column:TRANS_SUBCATA"`
	IndustryAddnInf      string    `gorm:"column:INDUSTRY_ADDN_INF"`
	TransCurrCd          string    `gorm:"column:TRANS_CURR_CD"`
	SecCtrlInf           string    `gorm:"column:SEC_CTRL_INF"`
	IcData               string    `gorm:"column:IC_DATA"`
	UdfFldPure           string    `gorm:"column:UDF_FLD_PURE"`
	CertifId             string    `gorm:"column:CERTIF_ID"`
	NetworkMgmtInfCd     string    `gorm:"column:NETWORK_MGMT_INF_CD"`
	OrigDataElemnt       string    `gorm:"column:ORIG_DATA_ELEMNT"`
	RcvInsIdCd           string    `gorm:"column:RCV_INS_ID_CD"`
	TfrInAcctNoPure      string    `gorm:"column:TFR_IN_ACCT_NO_PURE"`
	TfrInAcctTp          string    `gorm:"column:TFR_IN_ACCT_TP"`
	TfrOutAcctNoPure     string    `gorm:"column:TFR_OUT_ACCT_NO_PURE"`
	AcptInsResvPure      string    `gorm:"column:ACPT_INS_RESV_PURE"`
	TrrOutAcctTp         string    `gorm:"column:TRR_OUT_ACCT_TP"`
	IssInsIdCd           string    `gorm:"column:ISS_INS_ID_CD" downI:"11"`
	CardAttr             string    `gorm:"column:CARD_ATTR"`
	CardClass            string    `gorm:"column:CARD_CLASS" downI:"14"`
	CardMedia            string    `gorm:"column:CARD_MEDIA"`
	CardBin              string    `gorm:"column:CARD_BIN"`
	CardBrand            string    `gorm:"column:CARD_BRAND"`
	RoutInsIdCd          string    `gorm:"column:ROUT_INS_ID_CD"`
	AcptRegionCd         string    `gorm:"column:ACPT_REGION_CD"`
	BussRegionCd         string    `gorm:"column:BUSS_REGION_CD"`
	UsrNoTp              string    `gorm:"column:USR_NO_TP"`
	UsrNoRegionCd        string    `gorm:"column:USR_NO_REGION_CD"`
	UsrNoRegionAddnCd    string    `gorm:"column:USR_NO_REGION_ADDN_CD"`
	UsrNo                string    `gorm:"column:USR_NO"`
	SpInsIdCd            string    `gorm:"column:SP_INS_ID_CD"`
	IndustryInsIdCd      string    `gorm:"column:INDUSTRY_INS_ID_CD"`
	RoutIndustryInsIdCd  string    `gorm:"column:ROUT_INDUSTRY_INS_ID_CD"`
	IndustryMchntCd      string    `gorm:"column:INDUSTRY_MCHNT_CD"`
	IndustryTermCd       string    `gorm:"column:INDUSTRY_TERM_CD"`
	IndustryMchntTp      string    `gorm:"column:INDUSTRY_MCHNT_TP"`
	EntrustTp            string    `gorm:"column:ENTRUST_TP"`
	PmtMd                string    `gorm:"column:PMT_MD"`
	PmtTp                string    `gorm:"column:PMT_TP"`
	PmtNo                string    `gorm:"column:PMT_NO"`
	PmtMchntCd           string    `gorm:"column:PMT_MCHNT_CD"`
	PmtNoIndustryInsIdCd string    `gorm:"column:PMT_NO_INDUSTRY_INS_ID_CD"`
	PriAcctNoConv        string    `gorm:"column:PRI_ACCT_NO_CONV"`
	TransAtConv          string    `gorm:"column:TRANS_AT_CONV"`
	TransDtTmConv        string    `gorm:"column:TRANS_DT_TM_CONV"`
	TransSeqConv         string    `gorm:"column:TRANS_SEQ_CONV"`
	MchntTpConv          string    `gorm:"column:MCHNT_TP_CONV"`
	RetriRefNoConv       string    `gorm:"column:RETRI_REF_NO_CONV"`
	AcptInsIdCdConv      string    `gorm:"column:ACPT_INS_ID_CD_CONV"`
	TermIdConv           string    `gorm:"column:TERM_ID_CONV"`
	MchntCdConv          string    `gorm:"column:MCHNT_CD_CONV"`
	MchntNmConv          string    `gorm:"column:MCHNT_NM_CONV"`
	UdfFldPureConv       string    `gorm:"column:UDF_FLD_PURE_CONV"`
	SpInsIdCdConv        string    `gorm:"column:SP_INS_ID_CD_CONV"`
	ExpireDt             string    `gorm:"column:EXPIRE_DT"`
	SettleDt             string    `gorm:"column:SETTLE_DT"`
	TransFee             string    `gorm:"column:TRANS_FEE"`
	RespAuthId           string    `gorm:"column:RESP_AUTH_ID"`
	AcptRespCd           string    `gorm:"column:ACPT_RESP_CD"`
	AddnRespDataPure     string    `gorm:"column:ADDN_RESP_DATA_PURE"`
	AddnAtPure           string    `gorm:"column:ADDN_AT_PURE"`
	IssAddnDataPure      string    `gorm:"column:ISS_ADDN_DATA_PURE"`
	IcResDatCups         string    `gorm:"column:IC_RES_DAT_CUPS"`
	SwResvPure           string    `gorm:"column:SW_RESV_PURE"`
	IssInsResvPure       string    `gorm:"column:ISS_INS_RESV_PURE"`
	IndustryRespCd       string    `gorm:"column:INDUSTRY_RESP_CD"`
	DebtAt               string    `gorm:"column:DEBT_AT"`
	DtlInqData           string    `gorm:"column:DTL_INQ_DATA"`
	TransChnl            string    `gorm:"column:TRANS_CHNL"`
	InterchMdCd          string    `gorm:"column:INTERCH_MD_CD"`
	TransChkIn           string    `gorm:"column:TRANS_CHK_IN"`
	MchtStlmFlg          string    `gorm:"column:MCHT_STLM_FLG"`
	InsStlmFlg           string    `gorm:"column:INS_STLM_FLG"`
	MsgResvFld1          string    `gorm:"column:MSG_RESV_FLD1"`
	MsgResvFld2          string    `gorm:"column:MSG_RESV_FLD2"`
	MsgResvFld3          string    `gorm:"column:MSG_RESV_FLD3"`
	TransMth             int       `gorm:"column:TRANS_MTH"`
	RecUpdTs             time.Time `gorm:"column:REC_UPD_TS"`
	RecCrtTs             time.Time `gorm:"column:REC_CRT_TS"`
	ProdCd               string    `gorm:"column:PROD_CD" downI:"13"`
	TranTp               string    `gorm:"column:TRAN_TP"`
	BizCd                string    `gorm:"column:BIZ_CD"`
	RevelFlg             string    `gorm:"column:REVEL_FLG"`
	CancelFlg            string    `gorm:"column:CANCEL_FLG"`
	MsgResvFld4          string    `gorm:"column:MSG_RESV_FLD4"`
	MsgResvFld5          string    `gorm:"column:MSG_RESV_FLD5"`
	MsgResvFld6          string    `gorm:"column:MSG_RESV_FLD6"`
	MsgResvFld7          string    `gorm:"column:MSG_RESV_FLD7"`
	MsgResvFld8          string    `gorm:"column:MSG_RESV_FLD8"`
	MsgResvFld9          string    `gorm:"column:MSG_RESV_FLD9"`
}

func (t TfrTrnLog) TableName() string {
	return TfrTrnLogName
}

//DownloadInstitutionFile 根据开始时间跟结束时间查找交易流水
//@return Institution： 符合的实例
func (c ClearTxn) GetWithTime(db *gorm.DB, startTime, endTime string) ([]*ClearTxn, error) {
	clearTxn := []*ClearTxn{}
	err := db.Where("STLM_DATE <= ? AND STLM_DATE >= ?", endTime, startTime).Find(&clearTxn).Error
	if err != nil {
		return nil, err
	}
	return clearTxn, err
}

//查询商户一天的交易记录
//@params option 查询的limit跟offset， 第一个为limit， 第二个为page
//默认limit=10； page=0
//@params amountCond 可以为空， a < TRANS_AT AND b > TRANS_AT
func (t TfrTrnLog) GetWithLimit(db *gorm.DB, tfrTrnLog *TfrTrnLog, amountCond string, limit, page int64) ([]*TfrTrnLog, int64, string, error) {
	limit, offset := getLimitOffest(limit, page)
	logs := make([]*TfrTrnLog, 0)
	var count int64
	var total string

	err := db.Where(tfrTrnLog).Where(amountCond).Select(SelectInsTxnResp).Offset(offset).Limit(limit).Find(&logs).Error
	if err == gorm.ErrRecordNotFound {
		return logs, count, total, nil
	}
	if err != nil {
		return nil, count, total, err
	}

	rows, err := db.Table(tfrTrnLog.TableName()).Where(tfrTrnLog).Where(amountCond).Select("sum(TRANS_AT) as total, count(*) as count").Rows()
	if err != nil {
		return nil, count, total, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&total, &count)
		if err != nil {
			logrus.Error(err)
			return logs, 0, "", nil
		}
	}

	return logs, count, total, nil
}

func (t TfrTrnLog) Get(db *gorm.DB, tfrTrnLog *TfrTrnLog, amountCond string) ([]*TfrTrnLog, error) {
	logs := make([]*TfrTrnLog, 0)
	err := db.Where(tfrTrnLog).Where(amountCond).Select(SelectInsTxnResp).Find(&logs).Error
	if err == gorm.ErrRecordNotFound {
		return logs, nil
	}
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (t TfrTrnLog) GetByKeyRsp(db *gorm.DB, keyRsp string) (*TfrTrnLog, error) {
	tfrTrnLog := new(TfrTrnLog)
	err := db.Table(t.TableName()).Where("key_rsp = ?", keyRsp).First(tfrTrnLog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tfrTrnLog, nil
}

func getLimitOffest(limit, page int64) (int64, int64) {
	if limit == 0 {
		return 10, 0
	}

	offset := (page - 1) * limit
	return limit, offset
}
