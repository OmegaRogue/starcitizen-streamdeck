/*
Copyright © 2024 OmegaRogue <omegarogue@omegavoid.codes>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	sdk "github.com/OmegaRogue/streamdeck-sdk-go"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"starcitizen-streamdeck/internal/util"
	"starcitizen-streamdeck/pkg/goxdo"
	"starcitizen-streamdeck/pkg/sc"
)

var (
	cfgFile string
	rootCmd = NewRootCmd()
)

var data sc.Data

var funcMap = template.FuncMap{
	"localize": func(value string) string {
		loc, ok := data.Locale["english"][strings.Trim(strings.TrimSpace(value), "@")]
		if !ok {
			return value
		}
		return loc
	},
	"getBind": data.LookupBind,
	"empty": func(value string) bool {
		return strings.TrimSpace(value) == ""
	},
	"notEmpty": func(value string) bool {
		return strings.TrimSpace(value) != ""
	},
	"localizeKey": func(value string) string {
		return util.LocalizeKeyString(value, "de")
	},
}

//goland:noinspection GoUnusedExportedFunction
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "starcitizen-streamdeck",
		Short: "starcitizen-streamdeck",
		Long:  `starcitizen-streamdeck`,
		Run: func(cmd *cobra.Command, args []string) {

			data = sc.LoadData(viper.GetString("prefix"), sc.Version(viper.GetString("version")))

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create watcher")
			}
			defer util.DiscardErrorOnly(watcher.Close())

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						if event.Has(fsnotify.Write) {
							data.Rebinds = sc.LoadActionmap(viper.GetString("prefix"), sc.Version(viper.GetString("version")))
							RegenerateTemplates()
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						log.Err(err).Msg("error on watch")
					}
				}
			}()

			err = watcher.Add(sc.ActionMaps(viper.GetString("prefix"), sc.Version(viper.GetString("version"))))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to add actionmap to watcher")
			}

			RegenerateTemplates()

			err = sdk.Open()
			if err != nil {
				log.Fatal().Err(err).Msg("Error opening SDK Connection")
			}
			xdo := goxdo.NewXdo()
			defer xdo.Free()
			registerAction("static", func(action, context string, payload sdk.KeyPayload, deviceId string) {
				keyStr := sdk.JsonStringValue(payload.Settings, "function")
				key := data.LookupBind(strings.Trim(keyStr, "\""))
				xdo.SendKeysequenceWindowDown(goxdo.CURRENTWINDOW, util.FromScKey(key), 0)
			}, func(action, context string, payload sdk.KeyPayload, deviceId string) {
				keyStr := sdk.JsonStringValue(payload.Settings, "function")
				key := data.LookupBind(strings.Trim(keyStr, "\""))
				xdo.SendKeysequenceWindowUp(goxdo.CURRENTWINDOW, util.FromScKey(key), 0)
			})
			registerAction("macro", func(action, context string, payload sdk.KeyPayload, deviceId string) {
				keyStr := sdk.JsonStringValue(payload.Settings, "function")
				key := data.LookupBind(strings.Trim(keyStr, "\""))
				delay := payload.Settings.GetInt("delay")
				time.Sleep(time.Duration(delay) * time.Millisecond)
				xdo.SendKeysequenceWindow(goxdo.CURRENTWINDOW, util.FromScKey(key), 0)
			}, func(action, context string, payload sdk.KeyPayload, deviceId string) {})

			log.Trace().Msg("running")
			// Wait until the socket is closed, or SIGTERM/SIGINT is received
			sdk.Wait()
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weylus-desktop.yaml)")
	rootCmd.Flags().Int("port", 0, "Websocket port of the streamdeck software")
	rootCmd.Flags().String("pluginUUID", "", "UUID of the streamdeck plugin")
	rootCmd.Flags().String("registerEvent", "", "reserved")
	rootCmd.Flags().String("info", "", "reserved")

	return rootCmd
}

type Dict[K comparable, V any] map[K]V

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("exit with error")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		wd, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(wd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
		viper.SetDefault("version", sc.VersionLIVE)

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, err := fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			log.Err(err).Msg("print config file location")
		}
	}
}

func registerAction(id string, callback sdk.ActionHandler[sdk.KeyPayload], callbackUp sdk.ActionHandler[sdk.KeyPayload]) {
	sdk.RegisterActionDown("codes.omegavoid.starcitizen."+id, func(action, context string, payload sdk.KeyPayload, deviceId string) {
		callback(action, context, payload, deviceId)
		log.Info().Str("action", action).Str("direction", "down").Str("context", context).Int("state", payload.State).Bool("isInMultiAction", payload.IsInMultiAction).RawJSON("settings", []byte(payload.Settings.String())).Str("deviceId", deviceId).Msg("codes.omegavoid.starcitizen." + id)
	})
	sdk.RegisterActionUp("codes.omegavoid.starcitizen."+id, func(action, context string, payload sdk.KeyPayload, deviceId string) {
		callbackUp(action, context, payload, deviceId)
		log.Info().Str("action", action).Str("direction", "up").Str("context", context).Int("state", payload.State).Bool("isInMultiAction", payload.IsInMultiAction).RawJSON("settings", []byte(payload.Settings.String())).Str("deviceId", deviceId).Msg("codes.omegavoid.starcitizen." + id)
	})
}

func RegenerateTemplates() {
	templ, err := template.New("master").Funcs(funcMap).ParseGlob("*.gohtml")
	if err != nil {
		log.Err(err).Msg("parsing template")
	}

	templ.DefinedTemplates()
	templ.Templates()
	static, _ := os.OpenFile("PropertyInspector/StarCitizen/Static.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := templ.ExecuteTemplate(static, "static.gohtml", data.Profile); err != nil {
		log.Fatal().Err(err).Msg("static template failed to execute")
	}
	defer util.DiscardErrorOnly(static.Sync())
	defer util.DiscardErrorOnly(static.Close())
	macro, _ := os.OpenFile("PropertyInspector/StarCitizen/Macro.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := templ.ExecuteTemplate(macro, "macro.gohtml", data.Profile); err != nil {
		log.Fatal().Err(err).Msg("macro template failed to execute")
	}
	defer util.DiscardErrorOnly(macro.Sync())
	defer util.DiscardErrorOnly(macro.Close())
	dial, _ := os.OpenFile("PropertyInspector/StarCitizen/Dial.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err := templ.ExecuteTemplate(dial, "dial.gohtml", data.Profile); err != nil {
		log.Fatal().Err(err).Msg("dial template failed to execute")
	}
	defer util.DiscardErrorOnly(dial.Sync())
	defer util.DiscardErrorOnly(dial.Close())
}
