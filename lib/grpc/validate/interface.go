package grpc_validate

/*
	Интерфейс запроса
	Для выполнения валидации в запросе должен быть реализован метод Validate() error
*/
type ValidatableRequest interface {
	Validate() error
}
