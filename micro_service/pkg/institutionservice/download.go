package institutionservice

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	cleartxnM "userService/pkg/model/cleartxn"

	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
)

const filePath = "/tmp"
const defaultSheep = "Worksheet"
const fileZipPath = "/home/dj/zip"

var (
	titleTfrTrnLogs = "交易键值,商户编码,商户名称,交易日期,清算日期,交易时间,交易类型,收单机构代码,交易金额,交易卡号,发卡机构代码,终端编码,产品码,卡类型,交易状态,应答码"
	titleClearTxn   = "商户编码,交易日期,交易时间,清算日期,终端编号,交易类型,交易流水号,交易卡号,卡类型,交易本金,交易手续费,交易结算资金,应收差错费用,应付差费用,系统流水号,机构基准收入,机构实际收入,机构营销返佣,代理编码,会员号"
)

//通过日期将数据分为不同的文件
//@param fileDir test
//@return fileDir 下载完成之后存放临时文件的文件夹
func DownloadFileWithDay(clearTxn []*cleartxnM.ClearTxn) (string, error) {

	uid := uuid.NewV4().String()
	fileDirPath := path.Join(filePath, uid)
	os.MkdirAll(fileDirPath, os.ModePerm)

	if len(clearTxn) == 0 {
		return uid, nil
	}

	var dataIndex string

	for index := 0; index < len(clearTxn); index++ {
		dataIndex = clearTxn[index].StlmDate
		file := xlsx.NewFile()
		sheet, err := file.AddSheet(defaultSheep)
		if err != nil {
			return "", err
		}

		addTitle(sheet, titleClearTxn)
		var j = index
		for ; j < len(clearTxn); j++ {
			if dataIndex == clearTxn[j].StlmDate {
				row := sheet.AddRow()
				cell := row.AddCell()
				cell.Value = clearTxn[j].MchtCd
				cell = row.AddCell()
				cell.Value = clearTxn[j].StlmDate
				cell = row.AddCell()
				cell.Value = clearTxn[j].TransDateTime
				cell = row.AddCell()
				cell.Value = clearTxn[j].MaSettleDt.Format("20190430")
				cell = row.AddCell()
				cell.Value = clearTxn[j].TermID
				cell = row.AddCell()
				cell.Value = clearTxn[j].TxnDesc
				cell = row.AddCell()
				cell.Value = "交易流水号"
				cell = row.AddCell()
				cell.Value = "交易卡号"
				cell = row.AddCell()
				cell.Value = clearTxn[j].CardKindDis
				cell = row.AddCell()
				cell.Value = "交易本金"
				cell = row.AddCell()
				cell.Value = "交易手续费"
				cell = row.AddCell()
				cell.Value = "交易结算资金"
				cell = row.AddCell()
				cell.Value = "应收差错费用"
				cell = row.AddCell()
				cell.Value = "应付差费用"
				cell = row.AddCell()
				cell.Value = "系统流水号"
				cell = row.AddCell()
				cell.Value = "机构基准收入"
				cell = row.AddCell()
				cell.Value = "机构实际收入"
				cell = row.AddCell()
				cell.Value = "机构营销返佣"
				cell = row.AddCell()
				cell.Value = "代理编码"
				cell = row.AddCell()
				cell.Value = "会员号"
			} else {
				break
			}
		}
		index = j - 1

		filename := path.Join(fileDirPath, dataIndex+".xlsx")
		err = file.Save(filename)
		if err != nil {
			return "", err
		}
	}
	return uid, nil
}

func DownloadTfrTrnLogs(tfrTrnLog []*cleartxnM.TfrTrnLog) (string, error) {
	uid := uuid.NewV4().String()
	fileDirPath := path.Join(filePath, uid)
	os.MkdirAll(fileDirPath, os.ModePerm)

	if len(tfrTrnLog) == 0 {
		return uid, nil
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(defaultSheep)
	if err != nil {
		return "", err
	}

	addTitle(sheet, titleTfrTrnLogs)
	err = addBodyWithTfrTrnlogs(sheet, tfrTrnLog)
	if err != nil {
		return "", err
	}

	filename := path.Join(fileDirPath, "admin.xlsx")
	err = file.Save(filename)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func addBodyWithTfrTrnlogs(sheet *xlsx.Sheet, tfrTrnLog []*cleartxnM.TfrTrnLog) error {
	for i := range tfrTrnLog {
		var arrItems = [16]string{}
		trnType := reflect.TypeOf(*tfrTrnLog[i])
		trnValue := reflect.ValueOf(*tfrTrnLog[i])
		for j := 0; j < trnType.NumField(); j++ {
			tag, ok := trnType.Field(j).Tag.Lookup("downI")
			if !ok {
				continue
			}
			tagInt, err := strconv.ParseInt(tag, 10, 64)

			if err != nil {
				return err
			}

			val, ok := trnValue.Field(j).Interface().(string)
			if !ok {
				continue
			}
			arrItems[tagInt-1] = val
		}

		row := sheet.AddRow()
		for i := range arrItems {
			cell := row.AddCell()
			cell.Value = arrItems[i]
		}
	}
	return nil
}

func addTitle(sheet *xlsx.Sheet, title string) {
	s := strings.Split(title, ",")
	row := sheet.AddRow()
	for i, _ := range s {
		cell := row.AddCell()
		cell.Value = s[i]
	}
}

//@param dis 压缩前文件名称
//@param src 压缩后文件名称
func Compress(dis, src string) error {
	srcFile, err := os.Create(path.Join(fileZipPath, src))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	zw := zip.NewWriter(srcFile)
	defer zw.Close()

	disFileDir := path.Join(filePath, dis)
	filesInfo, err := ioutil.ReadDir(disFileDir)
	if err != nil {
		return errors.New("File not exist")
	}

	files := make([]string, 0)
	for _, f := range filesInfo {
		files = append(files, path.Join(disFileDir, f.Name()))
	}

	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
