/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: AGPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package logger

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/OmegaRogue/weylus-desktop/logger/journald"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupLogger() {
	var multi zerolog.LevelWriter
	journaldWriter := journald.NewBetterJournaldWriter()

	consoleWriter := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		FieldsExclude: []string{journald.ThreadFieldName},
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}
	multi = zerolog.MultiLevelWriter(consoleWriter, journaldWriter)
	log.Logger = log.Output(multi).With().Caller().Stack().Logger().Hook(journald.ThreadHook{})
	stdLogger := log.With().Str("component", "stdlog").Logger()
	stdlog.SetOutput(stdLogger)

	//zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.ErrorStackMarshaler = ESPackMarshalStack
}

func ESPackMarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	e, ok := err.(stackTracer)
	if !ok {
		return nil
	}

	pretty.Println(e.StackTrace())
	//It's mean when env=dev just print track
	if true {
		for _, frame := range e.StackTrace() {
			fmt.Printf("%+s:%d\r\n", frame, frame)
		}
	} else {
		return pkgerrors.MarshalStack(err)
	}
	return pkgerrors.MarshalStack(err)
}
