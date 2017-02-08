package storagemaster

import (
       "errors"
)

type Extras interface {
     Set(string, interface{}) error
     Get(string) (interface{}, error)
}

type Provider interface {
	Put(string, []byte, ...Extras) error
	Get(string) ([]byte, error)
	Delete(string) error
	Exists(string) (bool, error)
}

type StoragemasterExtras struct {
	Extras
	extras map[string]interface{}
}

func NewStoragemasterExtras() (*StoragemasterExtras, error){

     extras := make(map[string]interface{})
     
     e := StoragemasterExtras{
     	extras: extras,
     }

     return &e, nil
}

func (e *StoragemasterExtras) Get(key string) (interface{}, error){

     v, ok := e.extras[key]

     if !ok {
     	return nil, errors.New("Invalid key")
     }
     
     return v, nil
}

func (e *StoragemasterExtras) Set(key string, value interface{}) error{
     e.extras[key] = value
     return nil
}

