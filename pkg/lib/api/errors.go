package api

import "errors"

var InvalidPassword = errors.New("неправильный пароль")

var ErrCreateUser = errors.New("ошибка при создании пользователя")
