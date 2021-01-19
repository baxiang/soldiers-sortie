package staticservice

import (
	"reflect"
	"testing"
	"userService/pkg/model/static"
	"userService/pkg/pb"

	"golang.org/x/net/context"
)

func TestGetDicItemByCondition(t *testing.T) {
	cases := []struct {
		dicTypes []string
		dicNames []string
		dicCodes []string
		ret      []*static.DictionaryItem
	}{
		{
			dicTypes: []string{"TRANS_CD"},
			dicNames: []string{},
			dicCodes: []string{},
			ret: []*static.DictionaryItem{
				&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD1", DicCode: "100"},
				&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD2", DicCode: "101"},
				&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD3", DicCode: "102"},
			},
		},
		{
			dicTypes: []string{"PROD_CD"},
			dicNames: []string{"PROD_CD1"},
			dicCodes: []string{},
			ret: []*static.DictionaryItem{
				&static.DictionaryItem{DicType: "PROD_CD", DicName: "PROD_CD1", DicCode: "1000"},
			},
		},
		{
			dicTypes: []string{"BIZ_CD"},
			dicNames: []string{"BIZ_CD1"},
			dicCodes: []string{"10"},
			ret: []*static.DictionaryItem{
				&static.DictionaryItem{DicType: "BIZ_CD", DicName: "BIZ_CD1", DicCode: "10"},
			},
		},
	}

	for _, c := range cases {
		ret := getDicItemByCondition(c.dicTypes, c.dicNames, c.dicCodes)
		if !reflect.DeepEqual(ret, c.ret) {
			t.Errorf("want %v, got %v", c.ret, ret)
		}
	}
}

func TestGetBizCdByProdCd(t *testing.T) {
	cases := []struct {
		prodCd []string
		ret    []string
	}{
		{
			prodCd: []string{"1000"},
			ret:    []string{"10", "10", "11"},
		},
		{
			prodCd: []string{"10000"},
			ret:    []string{},
		},
	}

	for _, c := range cases {
		ret := getBizCdByProdCd(c.prodCd)
		if !reflect.DeepEqual(ret, c.ret) {
			t.Errorf("want %v, got %v", c.ret, ret)
		}
	}
}

func TestGetTransCdByProdAndBiz(t *testing.T) {
	cases := []struct {
		prodCd string
		bizCd  string
		ret    []string
	}{
		{
			prodCd: "1000",
			bizCd:  "10",
			ret:    []string{"100", "101"},
		},
		{
			prodCd: "10001",
			bizCd:  "10",
			ret:    []string{},
		},
	}

	for _, c := range cases {
		ret := getTransCdByProdAndBiz(c.prodCd, c.bizCd)
		if !reflect.DeepEqual(ret, c.ret) {
			t.Errorf("want %v, got %v", c.ret, ret)
		}
	}
}

func TestGetDicByInsCmpCd(t *testing.T) {
	cases := []struct {
		insCmpCd string
		ret      []string
	}{
		{
			insCmpCd: "1",
			ret:      []string{"2001", "2002"},
		},
		{
			insCmpCd: "2",
			ret:      []string{"2005"},
		},
		{
			insCmpCd: "0",
			ret:      []string{},
		},
	}

	for _, c := range cases {
		ret := getDicByInsCmpCd(c.insCmpCd)
		if !reflect.DeepEqual(ret, c.ret) {
			t.Errorf("want %v, got %v", c.ret, ret)
		}
	}
}

func TestIncludes(t *testing.T) {
	want := true
	have := includes([]string{"a", "b"}, "a")
	if want != have {
		t.Fatalf("want %t, got %t", want, have)
	}
	have = includes([]string{}, "")
	if want != have {
		t.Fatalf("want %t, got %t", want, have)
	}
	want = false
	have = includes([]string{"a"}, "c")
	if want != have {
		t.Fatalf("want %t, got %t", want, have)
	}
}

func TestMain(m *testing.M) {
	MyMap.dicItem = append(MyMap.dicItem, []*static.DictionaryItem{
		&static.DictionaryItem{DicType: "BIZ_CD", DicName: "BIZ_CD1", DicCode: "10"},
		&static.DictionaryItem{DicType: "BIZ_CD", DicName: "BIZ_CD2", DicCode: "11"},
		&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD1", DicCode: "100"},
		&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD2", DicCode: "101"},
		&static.DictionaryItem{DicType: "TRANS_CD", DicName: "TRANS_CD3", DicCode: "102"},
		&static.DictionaryItem{DicType: "PROD_CD", DicName: "PROD_CD1", DicCode: "1000"},
		&static.DictionaryItem{DicType: "PROD_CD", DicName: "PROD_CD2", DicCode: "1001"},
		&static.DictionaryItem{DicType: "INS_ID_CD", DicName: "INS_ID_CD1", DicCode: "2001"},
		&static.DictionaryItem{DicType: "INS_ID_CD", DicName: "INS_ID_CD2", DicCode: "2002"},
	}...)

	MyMap.prodBizTransMap = append(MyMap.prodBizTransMap, []*static.ProdBizTransMap{
		&static.ProdBizTransMap{ProdCd: "1000", BizCd: "10", TransCd: "100"},
		&static.ProdBizTransMap{ProdCd: "1000", BizCd: "10", TransCd: "101"},
		&static.ProdBizTransMap{ProdCd: "1000", BizCd: "11", TransCd: "102"},
	}...)

	MyMap.insInf = append(MyMap.insInf, []*static.InsInf{
		&static.InsInf{InsIDCd: "2001", InsCompanyCd: "1"},
		&static.InsInf{InsIDCd: "2002", InsCompanyCd: "1"},
		&static.InsInf{InsIDCd: "2005", InsCompanyCd: "2"},
	}...)

	MyMap.insProdBizFeeMapInf = append(MyMap.insProdBizFeeMapInf, []*static.InsProdBizFeeMapInf{
		&static.InsProdBizFeeMapInf{ProdCd: "1000", BizCd: "10"},
		&static.InsProdBizFeeMapInf{ProdCd: "1000", BizCd: "10"},
		&static.InsProdBizFeeMapInf{ProdCd: "1000", BizCd: "11"},
	}...)

	m.Run()
}

func TestGetDictionaryItem(t *testing.T) {
	TestSetService := NewSetService()
	ret, _ := TestSetService.GetDictionaryItem(context.Background(), &pb.StaticGetDictionaryItemReq{DicType: "BIZ_CD"})
	if len(ret.GetDictionaryItems) != 2 ||
		(*ret.GetDictionaryItems[0]).DicCode != "10" ||
		(*ret.GetDictionaryItems[1]).DicCode != "11" {
		t.Fatalf("Search result error")
	}

	ret, _ = TestSetService.GetDictionaryItem(context.Background(), &pb.StaticGetDictionaryItemReq{DicType: "BIZ_CD", DicCode: "10"})
	if len(ret.GetDictionaryItems) != 1 ||
		(*ret.GetDictionaryItems[0]).DicCode != "10" {
		t.Fatalf("Search result error")
	}

}

func TestGetDicByProdAndBiz(t *testing.T) {
	TestSetService := NewSetService()
	ret, _ := TestSetService.GetDicByProdAndBiz(context.Background(), &pb.StaticGetDicByProdAndBizReq{ProdCd: "1000"})
	if len(ret.GetDictionaryItems) != 2 ||
		(*ret.GetDictionaryItems[0]).DicCode != "10" ||
		(*ret.GetDictionaryItems[1]).DicCode != "11" ||
		(*ret.GetDictionaryItems[0]).DicType != "BIZ_CD" {
		t.Fatalf("Search result error")
	}

	ret, _ = TestSetService.GetDicByProdAndBiz(context.Background(), &pb.StaticGetDicByProdAndBizReq{ProdCd: "1000", BizCd: "13"})
	if len(ret.GetDictionaryItems) == 0 {
		t.Fatalf("Search result error")
	}

	ret, _ = TestSetService.GetDicByProdAndBiz(context.Background(), &pb.StaticGetDicByProdAndBizReq{ProdCd: "1000", BizCd: "10"})
	if len(ret.GetDictionaryItems) != 2 ||
		(*ret.GetDictionaryItems[0]).DicName != "TRANS_CD1" ||
		(*ret.GetDictionaryItems[1]).DicName != "TRANS_CD2" ||
		(*ret.GetDictionaryItems[0]).DicType != "TRANS_CD" {
		t.Fatalf("Search result error")
	}

}

func TestSGetDicByInsCmpCd(t *testing.T) {
	TestSetService := NewSetService()
	ret, _ := TestSetService.GetDicByInsCmpCd(context.Background(), &pb.StaticGetDicByInsCmpCdReq{InsCompanyCd: "1"})
	if len(ret.GetDictionaryItems) != 2 ||
		(*ret.GetDictionaryItems[0]).DicName != "INS_ID_CD1" ||
		(*ret.GetDictionaryItems[1]).DicName != "INS_ID_CD2" ||
		(*ret.GetDictionaryItems[0]).DicType != "INS_ID_CD" {
		t.Fatalf("Search result error")
	}

	ret, _ = TestSetService.GetDicByInsCmpCd(context.Background(), &pb.StaticGetDicByInsCmpCdReq{InsCompanyCd: "5"})
	if len(ret.GetDictionaryItems) != 0 {
		t.Fatalf("Search result error")
	}
}
