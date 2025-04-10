package api

import "errors"

var InvalidPassword = errors.New("неправильный пароль")

var ErrCreateUser = errors.New("ошибка при создании пользователя")

var ErrUserAlreadyExists = errors.New("пользователь с таким email уже существует")
