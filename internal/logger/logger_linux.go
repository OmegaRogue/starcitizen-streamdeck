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

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}
	log.Logger = log.Output(consoleWriter).With().Caller().Stack().Logger()
	stdLogger := log.With().Str("component", "stdlog").Logger()
	stdlog.SetOutput(stdLogger)

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

	if _, err := pretty.Println(e.StackTrace()); err != nil {
		return nil
	}
	// It's mean when env=dev just print track
	if true {
		for _, frame := range e.StackTrace() {
			fmt.Printf("%+s:%d\r\n", frame, frame)
		}
	} else {
		return pkgerrors.MarshalStack(err)
	}
	return pkgerrors.MarshalStack(err)
}
