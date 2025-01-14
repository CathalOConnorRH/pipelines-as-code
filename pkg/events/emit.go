package events

import (
	"context"
	"path/filepath"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewEventEmitter(client kubernetes.Interface, logger *zap.SugaredLogger) *EventEmitter {
	return &EventEmitter{
		client: client,
		logger: logger,
	}
}

type EventEmitter struct {
	client kubernetes.Interface
	logger *zap.SugaredLogger
}

func (e *EventEmitter) SetLogger(logger *zap.SugaredLogger) {
	e.logger = logger
}

func (e *EventEmitter) EmitMessage(repo *v1alpha1.Repository, loggerLevel zapcore.Level, message string) {
	if repo != nil {
		event := makeEvent(repo, loggerLevel, message)
		if _, err := e.client.CoreV1().Events(event.Namespace).Create(context.Background(), event, metav1.CreateOptions{}); err != nil {
			e.logger.Infof("Cannot create event: %s", err.Error())
		}
	}

	//nolint
	switch loggerLevel {
	case zapcore.DebugLevel:
		e.logger.Debug(message)
	case zapcore.ErrorLevel:
		e.logger.Error(message)
	case zapcore.InfoLevel:
		e.logger.Info(message)
	case zapcore.WarnLevel:
		e.logger.Warn(message)
	}
}

func makeEvent(repo *v1alpha1.Repository, loggerLevel zapcore.Level, message string) *v1.Event {
	event := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: repo.Name + "-",
			Namespace:    repo.Namespace,
			Labels: map[string]string{
				filepath.Join(pipelinesascode.GroupName, "repository"): repo.Name,
			},
		},
		Message: message,
		Type:    v1.EventTypeWarning,
		InvolvedObject: v1.ObjectReference{
			Kind:      "Repository",
			Name:      repo.Name,
			Namespace: repo.Namespace,
		},
	}
	if loggerLevel == zap.InfoLevel {
		event.Type = v1.EventTypeNormal
	}
	return event
}
