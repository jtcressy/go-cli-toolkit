// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Joel Cressy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package printers

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/samber/lo"
)

var _ ObjectPrinter = (*TablePrinter)(nil)

const (
	blue      = lipgloss.Color("021")
	purple    = lipgloss.Color("128")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("255")
)

var (
	DefaultTableHeaderStyle = lipgloss.NewStyle().
				Foreground(purple).
				Bold(true).
				Align(lipgloss.Center)
	DefaultTableCellWidth = 14
	DefaultTableCellStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Width(DefaultTableCellWidth).
				Align(lipgloss.Center).
				Foreground(lightGray)
	DefaultTableBorderType    = lipgloss.NormalBorder()
	DefaultTableBorderStyle   = lipgloss.NewStyle().Foreground(blue)
	DefaultTableCustomizeFunc = func(t *table.Table) *table.Table {
		return t.Border(DefaultTableBorderType).BorderStyle(DefaultTableBorderStyle)
	}
)

type TablePrinter struct {
	PrintOptions
	HeaderStyle        lipgloss.Style
	CellStyle          lipgloss.Style
	CellStyleFunc      func(style lipgloss.Style, row, col int, value string) lipgloss.Style
	TableCustomizeFunc func(t *table.Table) *table.Table
	TableReflectorFunc
}

// PrintObj implements ObjectPrinter.
func (p *TablePrinter) PrintObj(obj any, w io.Writer) (err error) {
	defer err2.Handle(&err, nil)

	headers, rows := try.To2(p.TableReflectorFunc(obj))

	if len(headers) == 0 || len(rows) == 0 {
		return nil
	}

	colWidths := lo.Reduce(rows, func(columnWidths []int, row []string, _ int) []int {
		return lo.Map(row, func(cell string, idx int) int {
			if len(cell) > columnWidths[idx] {
				return len(cell)
			}
			return columnWidths[idx]
		})
	}, lo.Map(headers, func(h string, _ int) int {
		return len(h)
	}))

	t := table.New().
		StyleFunc(func(row, col int) (style lipgloss.Style) {
			switch {
			case row == 0: // header
				return p.HeaderStyle.Width(colWidths[col] + style.GetHorizontalPadding())
			default:
				style = p.CellStyle
			}

			style = style.Width(colWidths[col] + style.GetHorizontalPadding())

			if p.CellStyleFunc != nil {
				return p.CellStyleFunc(
					style,
					row,
					col,
					rows[lo.If((row-1) < len(rows), row-1).Else(len(rows)-1)][col],
				)
			}
			return style
		})
	if p.TableCustomizeFunc != nil {
		t = p.TableCustomizeFunc(t)
	}
	if p.NoHeaders {
		t = t.Headers()
	} else {
		t = t.Headers(headers...)
	}
	_ = try.To1(io.WriteString(w, t.Data(table.NewStringData(rows...)).Render()))
	return nil
}

func NewTablePrinter(options PrintOptions) ObjectPrinter {
	printer := &TablePrinter{
		PrintOptions:       options,
		HeaderStyle:        DefaultTableHeaderStyle,
		CellStyle:          DefaultTableCellStyle,
		TableCustomizeFunc: DefaultTableCustomizeFunc,
		TableReflectorFunc: DefaultTableReflectorFunc,
	}
	return printer
}
