package receipt_generator

// func ConvertDataDoPertaminaToReceipt(data []*doPertaminaPB.ListModelData) []ListModelData {
// 	// Create a new array of interface{} and populate it with values from the struct array
// 	var listModelData []ListModelData
// 	for _, item := range data {
// 		var modelData []ModelData
// 		var tableData [][]string
// 		var columnWidth []uint32
// 		for _, mdata := range item.ModelData {
// 			modelData = append(modelData, ModelData{
// 				Key:            mdata.Key,
// 				Value:          mdata.Value,
// 				LenChar:        int(mdata.LenChar),
// 				Height:         int(mdata.Height),
// 				CutDot:         mdata.CutDot,
// 				IsTotalPayment: mdata.IsTotalPayment,
// 			})
// 		}

// 		if item.IsTable {
// 			for _, tdata := range item.TableSpec.TableData {
// 				var tableRow []string
// 				// for _, trow := range tdata.Data {
// 				tableRow = append(tableRow, tdata.Data...)
// 				// }
// 				tableData = append(tableData, tableRow)
// 			}

// 			for _, colWidth := range item.TableSpec.ColumnWidth {
// 				columnWidth = append(columnWidth, uint32(colWidth))
// 			}
// 		}

// 		listModelData = append(listModelData, ListModelData{
// 			HeaderData:  item.HeaderData,
// 			ModelData:   modelData,
// 			IsTable:     item.IsTable,
// 			TableData:   tableData,
// 			ColumnWidth: columnWidth,
// 		})
// 	}

// 	return listModelData
// }
