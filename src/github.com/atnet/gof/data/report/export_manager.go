/**
 * Copyright 2013 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package report

import (
	"database/sql"
	"encoding/json"
	"github.com/atnet/gof/db"
	"os"
	"strconv"
	"strings"
)

var baseExportPortal *DataExportPortal = &DataExportPortal{}

var _ IDataExportPortal = new(ExportItem)

// 导出项目
type ExportItem struct {
	Base      *DataExportPortal
	PortalKey string
	sqlConfig *ExportItemConfig
	//管理器
	ItemManager *ExportItemManager
}

func (portal *ExportItem) GetExportColumnIndexAndName(
	exportColumnNames []string) (dict map[string]string) {
	return portal.Base.GetExportColumnIndexAndName(exportColumnNames)
}

//获取统计数据
func (portal *ExportItem) GetTotalView(ht map[string]string) (row map[string]interface{}) {
	return nil
}

func (portal *ExportItem) GetSchemaAndData(ht map[string]string) (rows []map[string]interface{}, total int, err error) {
	total = 0
	var _rows *sql.Rows

	if portal.sqlConfig == nil {
		dir, _ := os.Getwd()
		portal.sqlConfig, err = LoadExportConfigFromXml(
			strings.Join([]string{dir, "/conf/query/", portal.PortalKey, ".xml"}, ""))
		if err != nil {
			portal.sqlConfig = nil
			return nil, 0, err
		}
	}

	_db := portal.ItemManager.DbGetter.GetDB()

	//初始化添加参数
	if _, e := ht["pageSize"]; !e {
		ht["pageSize"] = "100000"
	}
	if _, e := ht["pageIndex"]; !e {
		ht["pageIndex"] = "1"
	}

	pi, _ := ht["pageIndex"]
	ps, _ := ht["pageSize"]
	pageIndex, _ := strconv.Atoi(pi)
	pageSize, _ := strconv.Atoi(ps)

	if pageIndex > 0 {
		ht["page_start"] = strconv.Itoa((pageIndex - 1) * pageSize)
	} else {
		ht["page_start"] = "0"
	}
	ht["page_end"] = strconv.Itoa(pageIndex * pageSize)
	ht["page_size"] = strconv.Itoa(pageSize)

	//统计总行数
	if portal.sqlConfig.Total != "" {
		sql := SqlFormat(portal.sqlConfig.Total, ht)
		smt, err := _db.Prepare(sql)

		if err != nil {
			return nil, 0, err
		}

		row := smt.QueryRow()
		if row != nil {
			err = row.Scan(&total)
			if err != nil {
				return nil, total, err
			}
		}
	}

	//获得数据
	if portal.sqlConfig.Query != "" {
		sql := SqlFormat(portal.sqlConfig.Query, ht)

		smt, err := _db.Prepare(sql)

		if err != nil {
			return nil, total, err
		}
		_rows, err = smt.Query()
		if err != nil {
			return nil, total, err
		}
		defer _rows.Close()
	}

	return db.RowsToMarshalMap(_rows), total, err
}

//获取要导出的数据Json格式
func (portal *ExportItem) GetJsonData(ht map[string]string) string {
	result, err := json.Marshal(nil)
	if err != nil {
		return "{error:'" + err.Error() + "'}"
	}
	return string(result)
}

type IExportDbGetter interface {
	//获取数据库连接
	GetDB() *sql.DB
}

//导出项管理器
type ExportItemManager struct {
	//配置存放路径
	RootPath string
	//配置扩展名
	CfgFileExt string
	//数据库连接
	DbGetter IExportDbGetter //接口类型不需要加*
	//导出项集合
	exportItems map[string]ExportItem
}

//获取导出项
func (manager *ExportItemManager) GetExportItem(portalKey string) *ExportItem {

	if manager.RootPath == "" {
		manager.RootPath = "/conf/query/"
	}

	if manager.CfgFileExt == "" {
		manager.CfgFileExt = ".xml"
	}

	if itm, exist := manager.exportItems[portalKey]; exist {
		return &itm
	} else {
		dir, _ := os.Getwd()
		filePath := strings.Join(
			[]string{dir, manager.RootPath, portalKey, manager.CfgFileExt},
			"")
		if f, err := os.Stat(filePath); err == nil && f.IsDir() == false {
			itm := &ExportItem{
				PortalKey:   portalKey,
				Base:        baseExportPortal,
				ItemManager: manager}
			return itm
		}
	}
	return nil
}
