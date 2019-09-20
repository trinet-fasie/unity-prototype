package repository

import "gopkg.in/go-playground/validator.v8"

var validate = validator.New(&validator.Config{TagName: "validate"})
