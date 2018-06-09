package cherrylog

import (
	"github.com/containerum/cherry"
	"github.com/sirupsen/logrus"
	"strings"
)

// LogrusAdapter -- log origin and returning errors through logrus
type LogrusAdapter struct {
	*logrus.Entry
}

func (a *LogrusAdapter) Log(origin error, returning *cherry.Err) {
	a.WithError(origin).Errorf("returning %+v", returning)
}

// NewLogrusAdpater -- for more convenient usage
func NewLogrusAdapter(e *logrus.Entry) *LogrusAdapter {
	return &LogrusAdapter{Entry: e}
}

// Log(err).Errorf("unable to destroy Earth due Batman issues")
// Log(err, logger).Errorf("not enough energy units in the Death Star")
func Log(err *cherry.Err, optionalLogger ...logrus.FieldLogger) *logrus.Entry {
	var logger = func() logrus.FieldLogger {
		if len(optionalLogger) > 0 {
			return optionalLogger[0]
		}
		return logrus.StandardLogger()
	}()
	fields := make(map[string]interface{}, len(err.Fields))
	for k, v := range err.Fields {
		fields[k] = v
	}
	entry := logger.WithFields(fields).
		WithField("kind", err.ID.Kind).
		WithField("sid", err.ID.SID).
		WithField("status-http", err.StatusHTTP)
	if entry.Message == "" {
		entry.Message = err.Message + func() string {
			if len(err.Details) > 0 {
				return ": " + strings.Join(err.Details, ";")
			}
			return ""
		}()
	}
	return entry
}
