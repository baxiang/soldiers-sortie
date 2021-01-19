CREATE DATABASE if not exists apsmgm DEFAULT CHARACTER SET utf8mb4;

use apsmgm;

create table if not exists TBL_DICTIONARYITEM
(
    DIC_TYPE    VARCHAR(50)  not null,
    DIC_CODE    VARCHAR(50)  not null,
    DIC_NAME    VARCHAR(128) not null,
    DISP_ORDER  VARCHAR(50),
    UPDATE_TIME DATETIME,
    MEMO        VARCHAR(500),
    primary key (DIC_TYPE, DIC_CODE)
);

create table if not exists TBL_TFR_DYN_PARA
(
	PARA_NAME VARCHAR(100) not null,
	VALUE VARCHAR(10) not null,
	DESCRIPTION VARCHAR(1024),
    primary key(PARA_NAME)
);

create table if not exists TBL_TFR_TRN_LOG1
(
	TRANS_DT CHAR(8) not null,
	TRANS_MT CHAR(6) not null,
	SRC_QID INT not null,
	DES_QID INT,
	MA_TRANS_CD CHAR(10) not null,
	MA_TRANS_NM VARCHAR(60),
	KEY_RSP VARCHAR(48) not null,
	KEY_REVSAL VARCHAR(48),
	KEY_CANCEL VARCHAR(48),
	RESP_CD CHAR(2),
	TRANS_ST CHAR(10) not null,
	MA_TRANS_SEQ INT not null,
	ORIG_MA_TRANS_SEQ INT,
	ORIG_TRANS_SEQ CHAR(6),
	ORIG_TERM_SEQ CHAR(6),
	ORIG_TRANS_DT CHAR(8),
	MA_SETTLE_DT CHAR(8) not null,
	ACCESS_MD CHAR(4),
	MSG_TP CHAR(4),
	PRI_ACCT_NO CHAR(21),
	ACCT_TP CHAR(2),
	TRANS_PROC_CD CHAR(6),
	TRANS_AT CHAR(12),
	TRANS_TD_TM CHAR(10),
	TERM_SEQ CHAR(10),
	ACPT_TRANS_TM CHAR(6),
	ACPT_TRANS_DT CHAR(4),
	MCHNT_TP CHAR(4),
	POS_ENTRY_MD_CD CHAR(3),
	POS_COND_CD CHAR(2),
	ACPT_INS_ID_CD CHAR(11),
	FWD_INS_ID_CD CHAR(11),
	TERM_ID CHAR(8),
	MCHNT_CD CHAR(15),
	CARD_ACCPTR_NM NVARCHAR(40),
	RETRI_REF_NO CHAR(12),
	REQ_AUTH_ID CHAR(6),
	TRANS_SUBCATA VARCHAR(4),
	INDUSTRY_ADDN_INF VARCHAR(600),
	TRANS_CURR_CD CHAR(3),
	SEC_CTRL_INF CHAR(16),
	IC_DATA VARCHAR(103),
	UDF_FLD_PURE VARCHAR(100),
	CERTIF_ID VARCHAR(22),
	NETWORK_MGMT_INF_CD CHAR(3),
	ORIG_DATA_ELEMNT VARCHAR(42),
	RCV_INS_ID_CD CHAR(11),
	TFR_IN_ACCT_NO_PURE VARCHAR(28),
	TFR_IN_ACCT_TP CHAR(2),
	TFR_OUT_ACCT_NO_PURE VARCHAR(28),
	ACPT_INS_RESV_PURE VARCHAR(100),
	TRR_OUT_ACCT_TP CHAR(2),
	ISS_INS_ID_CD CHAR(11),
	CARD_ATTR CHAR(2),
	CARD_CLASS CHAR(2),
	CARD_MEDIA CHAR(2),
	CARD_BIN CHAR(14),
	CARD_BRAND CHAR(2),
	ROUT_INS_ID_CD CHAR(11),
	ACPT_REGION_CD CHAR(4),
	BUSS_REGION_CD CHAR(4),
	USR_NO_TP CHAR(2),
	USR_NO_REGION_CD CHAR(4),
	USR_NO_REGION_ADDN_CD CHAR(4),
	USR_NO VARCHAR(40),
	SP_INS_ID_CD CHAR(11),
	INDUSTRY_INS_ID_CD CHAR(11),
	ROUT_INDUSTRY_INS_ID_CD CHAR(11),
	INDUSTRY_MCHNT_CD CHAR(15),
	INDUSTRY_TERM_CD CHAR(8),
	INDUSTRY_MCHNT_TP CHAR(4),
	ENTRUST_TP CHAR,
	PMT_MD CHAR,
	PMT_TP CHAR(2),
	PMT_NO VARCHAR(40),
	PMT_MCHNT_CD CHAR(15),
	PMT_NO_INDUSTRY_INS_ID_CD CHAR(11),
	PRI_ACCT_NO_CONV CHAR(21),
	TRANS_AT_CONV CHAR(12),
	TRANS_DT_TM_CONV CHAR(10),
	TRANS_SEQ_CONV CHAR(6),
	MCHNT_TP_CONV CHAR(4),
	RETRI_REF_NO_CONV CHAR(12),
	ACPT_INS_ID_CD_CONV CHAR(11),
	TERM_ID_CONV CHAR(8),
	MCHNT_CD_CONV CHAR(15),
	MCHNT_NM_CONV VARCHAR(40),
	UDF_FLD_PURE_CONV VARCHAR(100),
	SP_INS_ID_CD_CONV CHAR(11),
	EXPIRE_DT CHAR(6),
	SETTLE_DT CHAR(8),
	TRANS_FEE CHAR(9),
	RESP_AUTH_ID CHAR(6),
	ACPT_RESP_CD CHAR(2),
	ADDN_RESP_DATA_PURE VARCHAR(25),
	ADDN_AT_PURE CHAR(40),
	ISS_ADDN_DATA_PURE VARCHAR(100),
	IC_RES_DAT_CUPS VARCHAR(103),
	SW_RESV_PURE VARCHAR(100),
	ISS_INS_RESV_PURE VARCHAR(100),
	INDUSTRY_RESP_CD VARCHAR(10),
	DEBT_AT VARCHAR(40),
	DTL_INQ_DATA VARCHAR(600),
	TRANS_CHNL CHAR(2),
	INTERCH_MD_CD CHAR,
	TRANS_CHK_IN CHAR,
	MCHT_STLM_FLG CHAR,
	INS_STLM_FLG CHAR,
	MSG_RESV_FLD1 VARCHAR(600),
	MSG_RESV_FLD2 VARCHAR(600),
	MSG_RESV_FLD3 VARCHAR(600),
	TRANS_MTH INT not null,
	REC_UPD_TS DATETIME,
	REC_CRT_TS DATETIME,
	PROD_CD VARCHAR(10),
	TRAN_TP VARCHAR(10),
	BIZ_CD VARCHAR(10),
	REVEL_FLG VARCHAR(10),
	CANCEL_FLG VARCHAR(10),
	MSG_RESV_FLD4 VARCHAR(100),
	MSG_RESV_FLD5 VARCHAR(100),
	MSG_RESV_FLD6 VARCHAR(100),
	MSG_RESV_FLD7 VARCHAR(100),
	MSG_RESV_FLD8 VARCHAR(100),
	MSG_RESV_FLD9 VARCHAR(100),
    primary key (KEY_RSP)
)


