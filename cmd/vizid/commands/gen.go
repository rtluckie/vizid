package commands

import (
	"fmt"

	"github.com/ryanl/vizid/internal/generator"
	"github.com/ryanl/vizid/internal/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	compYear   bool
	compMonth  bool
	compDay    bool
	compHour   bool
	compMinute bool
	compSecond bool
	compMs     bool
	compUUID   bool
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new VIZID (visual form, suitable for filenames)",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := model.Options{
			Timezone: viper.GetString("timezone"),
			Warn:     viper.GetBool("warn"),
			Custom:   viper.GetBool("custom"),
			Components: model.Components{
				Year:   viper.GetBool("components.year"),
				Month:  viper.GetBool("components.month"),
				Day:    viper.GetBool("components.day"),
				Hour:   viper.GetBool("components.hour"),
				Minute: viper.GetBool("components.minute"),
				Second: viper.GetBool("components.second"),
				Ms:     viper.GetBool("components.ms"),
				UUID:   viper.GetBool("components.uuid"),
			},
		}

		// If custom mode is enabled, override from flags.
		if opts.Custom {
			opts.Components.Year = compYear
			opts.Components.Month = compMonth
			opts.Components.Day = compDay
			opts.Components.Hour = compHour
			opts.Components.Minute = compMinute
			opts.Components.Second = compSecond
			opts.Components.Ms = compMs
			opts.Components.UUID = compUUID
		}

		id, ascii, warnMsg, err := generator.Generate(opts)
		if err != nil {
			return err
		}
		if opts.Warn && warnMsg != "" {
			fmt.Println("WARN:", warnMsg)
		}
		fmt.Println(id)
		fmt.Println(ascii)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().BoolVar(&compYear, "year", true, "include year")
	genCmd.Flags().BoolVar(&compMonth, "month", true, "include month")
	genCmd.Flags().BoolVar(&compDay, "day", true, "include day")
	genCmd.Flags().BoolVar(&compHour, "hour", true, "include hour")
	genCmd.Flags().BoolVar(&compMinute, "minute", true, "include minute")
	genCmd.Flags().BoolVar(&compSecond, "second", true, "include second")
	genCmd.Flags().BoolVar(&compMs, "ms", true, "include milliseconds")
	genCmd.Flags().BoolVar(&compUUID, "uuid", true, "include uuid")

	_ = viper.BindPFlag("components.year", genCmd.Flags().Lookup("year"))
	_ = viper.BindPFlag("components.month", genCmd.Flags().Lookup("month"))
	_ = viper.BindPFlag("components.day", genCmd.Flags().Lookup("day"))
	_ = viper.BindPFlag("components.hour", genCmd.Flags().Lookup("hour"))
	_ = viper.BindPFlag("components.minute", genCmd.Flags().Lookup("minute"))
	_ = viper.BindPFlag("components.second", genCmd.Flags().Lookup("second"))
	_ = viper.BindPFlag("components.ms", genCmd.Flags().Lookup("ms"))
	_ = viper.BindPFlag("components.uuid", genCmd.Flags().Lookup("uuid"))
}
