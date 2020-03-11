package printer

import (
	"tabtoy/v2/i18n"
	"tabtoy/v2/model"
	"tabtoy/util"
)

func valueWrapperPbt(t model.FieldType, node *model.Node) string {

	switch t {
	case model.FieldType_String:
		return util.StringEscape(node.Value)
	}

	return node.Value
}

type pbtPrinter struct {
}

func (self *pbtPrinter) Run(g *Globals) *Stream {

	bf := NewStream()
	bf.Printf("# Generated by github.com/davyxu/tabtoy\n")
	bf.Printf("# Version: %s\n", g.Version)

	for _, tab := range g.Tables {

		if !tab.LocalFD.MatchTag(".pbt") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		if !printTablePBT(bf, tab) {
			return nil
		}
	}

	return bf
}

func printTablePBT(bf *Stream, tab *model.Table) bool {

	if len(tab.Recs) == 0 {
		return true
	}

	bf.Printf("%s: [\n", tab.LocalFD.Name)

	// 遍历每一行
	for recIndex, r := range tab.Recs {

		bf.Printf("	{")

		// 遍历每一列
		for rootFieldIndex, node := range r.Nodes {

			if node.IsRepeated {
				bf.Printf("%s:[ ", node.Name)
			} else {
				bf.Printf("%s: ", node.Name)
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				if node.IsRepeated {

					// repeated 值序列
					for arrIndex, valueNode := range node.Child {

						bf.Printf("%s", valueWrapperPbt(node.Type, valueNode))

						// 多个值分割
						if arrIndex < len(node.Child)-1 {
							bf.Printf(", ")
						}

					}
				} else {
					// 单值
					valueNode := node.Child[0]

					if !node.SugguestIgnore {
						bf.Printf("%s", valueWrapperPbt(node.Type, valueNode))
					}

				}

			} else {

				// 遍历repeated的结构体
				for structIndex, structNode := range node.Child {

					// 结构体开始
					bf.Printf("{ ")

					// 遍历一个结构体的字段
					for structFieldIndex, fieldNode := range structNode.Child {

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						bf.Printf("%s: %s", fieldNode.Name, valueWrapperPbt(fieldNode.Type, valueNode))

						// 结构体字段分割
						if structFieldIndex < len(structNode.Child)-1 {
							bf.Printf(", ")
						}

					}

					// 结构体结束
					bf.Printf(" }")

					// 多个结构体分割
					if structIndex < len(node.Child)-1 {
						bf.Printf(", ")
					}

				}

			}

			if node.IsRepeated {
				bf.Printf(" ]")
			}

			// 根字段分割
			if rootFieldIndex < len(r.Nodes)-1 {
				bf.Printf(", ")
			}

		}

		bf.Printf(" }")

		// 根字段分割
		if recIndex < len(tab.Recs)-1 {
			bf.Printf(", ")
		}

		bf.Printf("\n")

	}

	bf.Printf("]\n\n")

	return true

}

func init() {

	RegisterPrinter("pbt", &pbtPrinter{})

}
