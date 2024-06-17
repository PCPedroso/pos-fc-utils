package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1          TestEvent
	event2          TestEvent
	handler1        TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler1 = TestEventHandler{ID: 1}
	suite.handler1 = TestEventHandler{ID: 1}
	suite.handler3 = TestEventHandler{ID: 3}

	suite.event1 = TestEvent{Name: "test1", Payload: "teste1"}
	suite.event2 = TestEvent{Name: "test2", Payload: "teste2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcherRegister() {
	// Registrando o handler1
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	// Verificando a quantidade de eventos registrados
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	// Registrando o handler2
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	assert.Equal(suite.T(), &suite.handler1, suite.eventDispatcher.handlers[suite.event1.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][1])
}

// Testando eventos registrados em deplicidade
func (suite *EventDispatcherTestSuite) TestEventDispatcherRegisterWithSameHandler() {
	// Registra o handler1
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	// Registra novamente o handler1 e deve retornar o erro de evento já registrado
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Equal(ErrHandlerAlreadyRegistred, err)
	// Ao final não deve registrar este segundo evento, e conter apenas um evento registrado
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
}

// Limpando os Eventos
func (suite *EventDispatcherTestSuite) TestEventDispatcherClear() {
	// Event1
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	// Event2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler1)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
