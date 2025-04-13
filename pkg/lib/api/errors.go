package api

import "errors"

var InvalidPassword = errors.New("неправильный пароль")

var ErrCreateUser = errors.New("ошибка при создании пользователя")

var ErrUserAlreadyExists = errors.New("пользователь с таким email уже существует")

var ErrorCreateQueryString = errors.New("ошибка при создании запроса строки базы данных")

var ErrCreatePvz = errors.New("ошибка при создании пункта выдачи заказов")

var (
	ErrNotFoundProduct    = errors.New("продукта с тами id нету")
	ErrNotFoundAcceptance = errors.New("пункта выдачи с таким id нету")
)
