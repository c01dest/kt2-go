package main

import (
	"fmt"
	"strings"
)

func main() {
	demoFacade()
	demoBuilder()
	demoVisitor()
	demoCommand()
	demoChainOfResponsibility()
	demoFactoryMethod()
	demoStrategy()
	demoState()
}

// --- 1. ФАСАД ---

type Inventory struct{}

func (i *Inventory) CheckStock(item string) bool {
	fmt.Println("[Склад] Проверка наличия товара:", item)
	return true
}

type PaymentGateway struct{}

func (p *PaymentGateway) Pay(amount float64) bool {
	fmt.Printf("[Оплата] Списано %.2f руб.\n", amount)
	return true
}

type OrderNotification struct{}

func (n *OrderNotification) SendSMS() {
	fmt.Println("[Уведомления] Отправка SMS с подтверждением заказа")
}

type OrderFacade struct {
	inventory *Inventory
	payment   *PaymentGateway
	notify    *OrderNotification
}

func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		inventory: &Inventory{},
		payment:   &PaymentGateway{},
		notify:    &OrderNotification{},
	}
}

func (o *OrderFacade) PlaceOrder(item string, amount float64) {
	if o.inventory.CheckStock(item) && o.payment.Pay(amount) {
		o.notify.SendSMS()
		fmt.Println("Результат: Заказ успешно оформлен!")
	}
}

func demoFacade() {
	fmt.Println("--- 1. Фасад ---")
	order := NewOrderFacade()
	order.PlaceOrder("Игровой ноутбук", 120000)
}

// --- 2. СТРОИТЕЛЬ ---

type Server struct {
	Host string
	Port int
	TLS  bool
}

type ServerBuilder struct {
	server Server
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{server: Server{Host: "localhost", Port: 8080, TLS: false}}
}

func (b *ServerBuilder) SetHost(host string) *ServerBuilder {
	b.server.Host = host
	return b
}

func (b *ServerBuilder) SetPort(port int) *ServerBuilder {
	b.server.Port = port
	return b
}

func (b *ServerBuilder) EnableTLS() *ServerBuilder {
	b.server.TLS = true
	return b
}

func (b *ServerBuilder) Build() Server {
	return b.server
}

func demoBuilder() {
	fmt.Println("\n--- 2. Строитель ---")
	server := NewServerBuilder().
		SetHost("127.0.0.1").
		SetPort(443).
		EnableTLS().
		Build()

	fmt.Printf("Конфигурация сервера: %+v\n", server)
}

// --- 3. ПОСЕТИТЕЛЬ ---

type Shape interface {
	Accept(Visitor)
}

type Square struct{ Side float64 }

func (s *Square) Accept(v Visitor) { v.VisitSquare(s) }

type Circle struct{ Radius float64 }

func (c *Circle) Accept(v Visitor) { v.VisitCircle(c) }

type Visitor interface {
	VisitSquare(*Square)
	VisitCircle(*Circle)
}

type AreaCalculator struct {
	Area float64
}

func (a *AreaCalculator) VisitSquare(s *Square) {
	a.Area = s.Side * s.Side
}

func (a *AreaCalculator) VisitCircle(c *Circle) {
	a.Area = 3.14 * c.Radius * c.Radius
}

func demoVisitor() {
	fmt.Println("\n--- 3. Посетитель ---")
	shapes := []Shape{&Square{Side: 4}, &Circle{Radius: 3}}
	calc := &AreaCalculator{}

	for _, shape := range shapes {
		shape.Accept(calc)
		fmt.Printf("Рассчитанная площадь: %.2f\n", calc.Area)
	}
}

// --- 4. КОМАНДА ---

type Command interface {
	Execute()
}

type Light struct{}

func (l *Light) On()  { fmt.Println("Свет включен") }
func (l *Light) Off() { fmt.Println("Свет выключен") }

type LightOnCommand struct{ light *Light }

func (c *LightOnCommand) Execute() { c.light.On() }

type LightOffCommand struct{ light *Light }

func (c *LightOffCommand) Execute() { c.light.Off() }

type RemoteControl struct {
	command Command
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

func demoCommand() {
	fmt.Println("\n--- 4. Команда ---")
	light := &Light{}
	onCmd := &LightOnCommand{light: light}
	offCmd := &LightOffCommand{light: light}

	remote := &RemoteControl{command: onCmd}
	remote.PressButton()

	remote.command = offCmd
	remote.PressButton()
}

// --- 5. ЦЕПОЧКА ВЫЗОВОВ ---

type Handler interface {
	SetNext(Handler) Handler
	Execute(*Request)
}

type Request struct {
	Authenticated bool
	Authorized    bool
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(next Handler) Handler {
	h.next = next
	return next
}

type AuthHandler struct{ BaseHandler }

func (h *AuthHandler) Execute(r *Request) {
	if !r.Authenticated {
		fmt.Println("[Цепочка] Ошибка: Пользователь не вошел в систему")
		return
	}
	fmt.Println("[Цепочка] 1. Успешная аутентификация")
	if h.next != nil {
		h.next.Execute(r)
	}
}

type RoleHandler struct{ BaseHandler }

func (h *RoleHandler) Execute(r *Request) {
	if !r.Authorized {
		fmt.Println("[Цепочка] Ошибка: Отказано в доступе")
		return
	}
	fmt.Println("[Цепочка] 2. Успешная авторизация (проверка прав)")
	if h.next != nil {
		h.next.Execute(r)
	}
}

func demoChainOfResponsibility() {
	fmt.Println("\n--- 5. Цепочка вызовов ---")
	auth := &AuthHandler{}
	role := &RoleHandler{}
	auth.SetNext(role)

	req := &Request{Authenticated: true, Authorized: true}
	auth.Execute(req)
}

// --- 6. ФАБРИЧНЫЙ МЕТОД ---

type Notifier interface {
	Send(msg string)
}

type EmailNotification struct{}

func (e *EmailNotification) Send(msg string) {
	fmt.Println("[Email] Отправлено сообщение:", msg)
}

type SMSNotification struct{}

func (s *SMSNotification) Send(msg string) {
	fmt.Println("[SMS] Отправлено сообщение:", msg)
}

func NewNotification(channel string) (Notifier, error) {
	switch strings.ToLower(channel) {
	case "email":
		return &EmailNotification{}, nil
	case "sms":
		return &SMSNotification{}, nil
	default:
		return nil, fmt.Errorf("неизвестный канал: %s", channel)
	}
}

func demoFactoryMethod() {
	fmt.Println("\n--- 6. Фабричный метод ---")
	notifier, err := NewNotification("sms")
	if err == nil {
		notifier.Send("Код доступа: 9982")
	}
}

// --- 7. СТРАТЕГИЯ ---

type PaymentStrategy interface {
	Pay(amount int)
}

type CreditCardStrategy struct{}

func (c *CreditCardStrategy) Pay(amount int) {
	fmt.Printf("[Оплата] Оплачено %d руб. банковской картой\n", amount)
}

type CryptoStrategy struct{}

func (k *CryptoStrategy) Pay(amount int) {
	fmt.Printf("[Оплата] Оплачено %d руб. криптовалютой USDT\n", amount)
}

type ShoppingCart struct {
	strategy PaymentStrategy
}

func (s *ShoppingCart) SetStrategy(strategy PaymentStrategy) {
	s.strategy = strategy
}

func (s *ShoppingCart) Checkout(amount int) {
	s.strategy.Pay(amount)
}

func demoStrategy() {
	fmt.Println("\n--- 7. Стратегия ---")
	cart := &ShoppingCart{}

	cart.SetStrategy(&CreditCardStrategy{})
	cart.Checkout(2500)

	cart.SetStrategy(&CryptoStrategy{})
	cart.Checkout(5000)
}

// --- 8. СОСТОЯНИЕ ---

type State interface {
	PressButton(*PlayerContext)
}

type PlayerContext struct {
	state State
}

func (c *PlayerContext) SetState(s State) {
	c.state = s
}

func (c *PlayerContext) PressButton() {
	c.state.PressButton(c)
}

type PlayingState struct{}

func (p *PlayingState) PressButton(c *PlayerContext) {
	fmt.Println("[Плеер] Нажата кнопка -> На паузу")
	c.SetState(&PausedState{})
}

type PausedState struct{}

func (p *PausedState) PressButton(c *PlayerContext) {
	fmt.Println("[Плеер] Нажата кнопка -> Запуск воспроизведения")
	c.SetState(&PlayingState{})
}

func demoState() {
	fmt.Println("\n--- 8. Состояние ---")
	player := &PlayerContext{state: &PausedState{}}

	player.PressButton()
	player.PressButton()
}
