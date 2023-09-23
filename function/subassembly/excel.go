package subassembly

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	log "main/function/logger"
	"time"
)

// MakeApiResultDown 用于对API结果下载到Excel里面
func MakeApiResultDown(name string, data []string) {
	// 指定要写入的 Excel 文件名
	excelName := name
	// 要写入的字符串数组
	datas := data
	day := time.Now().Format("2006.1.2")
	fileName := fmt.Sprintf("./result/%v-API测绘结果下载-%v.xlsx", excelName, day)

	// 打开 Excel 文件，如果文件不存在则创建新文件
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		file = xlsx.NewFile()
	}

	// 获取或创建一个工作表
	sheet, found := file.Sheet["Sheet1"]
	if !found {
		sheet, err = file.AddSheet("Sheet1")
		if err != nil {
			log.Debug(fmt.Sprintf("创建工作表失败: %s\n", err))
			return
		}
	}

	// 创建一行，并将数据写入单元格
	for _, value := range datas {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = value
	}

	// 保存 Excel 文件
	err = file.Save(fileName)
	if err != nil {
		log.Debug(fmt.Sprintf("保存 Excel 文件失败: %s", err))
		return
	}

	log.Debug(fmt.Sprintf("Excel 文件已追加保存：%s", fileName))
}
