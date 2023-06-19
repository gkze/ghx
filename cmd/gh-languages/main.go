package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"ghx/languages"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/exp/slices"
)

var (
	headers []string = []string{
		"ID", "NAME", "TYPE", "COLOR", "EXTENSIONS", "TEXTMATE_SCOPE",
		"ACE_MODE",
	}
	rootCmd   *cobra.Command
	sortField *string
)

func writeTabRow(tw *tabwriter.Writer, row []string, maxColWidth int) (int, error) {
	for idx, elem := range row {
		if len(elem) >= maxColWidth {
			row[idx] = elem[0:maxColWidth] + "..."
		}
	}

	return tw.Write([]byte(strings.Join(row, "\t") + "\n"))
}

func rootRunE(cmd *cobra.Command, args []string) error {
	if !slices.Contains(headers, strings.ToUpper(*sortField)) {
		return fmt.Errorf("%s is not a valid sort key", *sortField)
	}

	maxwidth, _, err := terminal.GetSize(0)
	if err != nil {
		return err
	}

	maxColWidth := maxwidth / len(headers)

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	if _, err := writeTabRow(tw, headers, maxColWidth); err != nil {
		return err
	}

	allLanguages, err := languages.GetLanguages(
		languages.GetLanguagesOpts{SortField: strings.ToUpper(*sortField)},
	)
	if err != nil {
		return err
	}

	for _, lang := range allLanguages {
		if _, err := writeTabRow(tw, []string{
			strconv.FormatUint(lang.ID, 10),
			lang.Name, string(lang.Type),
			languages.RGBA64ToHexString(lang.Color),
			strings.Join(lang.Extensions, ", "), lang.TextmateScope,
			lang.AceMode,
		}, maxColWidth); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func init() {
	rootCmd = &cobra.Command{
		Use: "gh-languages", Short: "Languages", RunE: rootRunE,
	}
	sortField = rootCmd.PersistentFlags().StringP("sort", "s", "ID",
		fmt.Sprintf(
			"Field to sort on. Options:\n - %s\n",
			strings.Join(headers, "\n - "),
		),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
