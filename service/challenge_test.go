package service

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCalculateStatus(t *testing.T) {

	Convey("Verify calculate result", t, func() {

		Convey("Test wrong file path", func() {
			err := CalculateStatus("/tmp/data.json")
			So(strings.Contains(err.Error(), "failed to locate data file"), ShouldBeTrue)
		})

		Convey("Test wrong file format", func() {
			err := CalculateStatus("../data/test_error_data.json")
			So(strings.Contains(err.Error(), "failed to parse data file"), ShouldBeTrue)
		})

		Convey("Test zero revenue", func() {

			old := os.Stdout // keep backup of the real stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			CalculateStatus("../data/test_no_revenue.json")
			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = old // restoring the real stdout
			out := <-outC
			So(strings.Contains(out, "Gross Profit Margin: 0%"), ShouldBeTrue)
			So(strings.Contains(out, "Net Profit Margin: 0%"), ShouldBeTrue)
		})

		Convey("Test zero liabilities", func() {

			old := os.Stdout // keep backup of the real stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			CalculateStatus("../data/test_no_liabilities.json")
			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = old // restoring the real stdout
			out := <-outC
			So(strings.Contains(out, "Working Capital Ratio: 0%"), ShouldBeTrue)
		})

		Convey("Test standard data", func() {

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			CalculateStatus("../data/data.json")
			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = old
			out := <-outC
			So(strings.Contains(out, "Revenue: $32,431"), ShouldBeTrue)
			So(strings.Contains(out, "Expense: $36,530"), ShouldBeTrue)
			So(strings.Contains(out, "Gross Profit Margin: 0%"), ShouldBeTrue)
			So(strings.Contains(out, "Net Profit Margin: -13%"), ShouldBeTrue)
			So(strings.Contains(out, "Working Capital Ratio: 119%"), ShouldBeTrue)
		})

	})
}
