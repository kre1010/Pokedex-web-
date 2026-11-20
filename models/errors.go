package models 	

import (
    "errors"
)

var (
       ErrNotFound  = errors.New("élément introuvable")
	   ErrInvalidInput = errors.New("données invalides") 
       ErrTimeout = errors.New("temps de chargement dépasser")
	   ErrUnavailable = errors.New("service pas accésible")
       ErrBadRequest = errors.New("mauvaise réponse")
)