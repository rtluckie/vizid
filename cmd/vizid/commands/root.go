package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	tzFlag   string
	warnFlag bool
	custom   bool
)

var rootCmd = &cobra.Command{
	Use:   "vizid",
	Short: "VIZID: visual, sortable timestamps + IDs (Unicode filenames)",
	Long:  "Generate and decode VIZIDs (visual IDs) designed to sort correctly in filenames.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.config/vizid/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&tzFlag, "timezone", "tz", "UTC", "timezone (IANA name like America/Chicago, UTC offset like +02:00, or UTC)")
	rootCmd.PersistentFlags().BoolVar(&warnFlag, "warn", true, "warn if sort order might be broken by custom component selection")
	rootCmd.PersistentFlags().BoolVarP(&custom, "custom", "C", false, "enable custom component selection flags")

	_ = viper.BindPFlag("timezone", rootCmd.PersistentFlags().Lookup("timezone"))
	_ = viper.BindPFlag("warn", rootCmd.PersistentFlags().Lookup("warn"))
	_ = viper.BindPFlag("custom", rootCmd.PersistentFlags().Lookup("custom"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home + "/.config/vizid")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetDefault("timezone", "UTC")
	viper.SetDefault("warn", true)
	viper.SetDefault("custom", false)

	// Components defaults
	viper.SetDefault("components.year", true)
	viper.SetDefault("components.month", true)
	viper.SetDefault("components.day", true)
	viper.SetDefault("components.hour", true)
	viper.SetDefault("components.minute", true)
	viper.SetDefault("components.second", true)
	viper.SetDefault("components.ms", true)
	viper.SetDefault("components.uuid", true)

	if err := viper.ReadInConfig(); err == nil {
		// config loaded
	}
}
