package cache

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/lib/contx"
)

/*
PutIntoCache puts a value at system's cache
*/
func PutIntoCache(key string, value interface{}) {
	s := fmt.Sprintf("%.0f", (time.Hour * 48).Seconds())
	timeout, _ := strconv.Atoi(s)

	err := PutIntoCacheWithTimeout(key, timeout, value)
	if err != nil {
		log.Printf("[PutIntoCache] Error putting into cache: [%s]. Key: [%s] - Value: [%v]\n", err.Error(), key, value)
        return
	}
	log.Printf("[PutIntoCache] Cache inputed. Key: [%s] - Value: [%v]\n", key, value)
}

/*
PutIntoCache puts a value at system's cache
*/
func PutIntoCacheWithTimeout(key string, timeoutInSeconds int, value interface{}) error {
	ctx := contx.GetContext()
	cache := ctx.Cache

	err := cache.Put(key, value, int64(timeoutInSeconds))
	if err != nil {
		log.Printf("[PutIntoCacheWithTimeout] Error putting into cache: [%s]. Key: [%s] - Value: [%v] - Timeout: [%d]\n", err.Error(), key, value, timeoutInSeconds)
		return err
	}
	log.Printf("[PutIntoCacheWithTimeout] Cache inputed. Key: [%s] - Value: [%v] - Timeout: [%d]\n", key, value, timeoutInSeconds)
	return nil
}

/*
RemoveFromCache removes an item from cache
*/
func RemoveFromCache(key string) {
	ctx := contx.GetContext()
	cache := ctx.Cache
	if cache.IsExist(key) {
		cache.Delete(key)
	}
}

/*
GetValueFromCache gets a value from the system's cache
*/
func GetValueFromCache(key string) interface{} {
	ctx := contx.GetContext()
	cache := ctx.Cache
	if cache.IsExist(key) {
		return cache.Get(key)
	}
	return nil
}

/*
GetIntValueFromCache gets an int value from the system's cache
*/
func GetIntValueFromCache(key string) (retorno int) {
	a := GetValueFromCache(key)
	if a != nil {
		retorno = a.(int)
	}
	return
}

/*
GetFloatValueFromCache gets an float value from the system's cache
*/
func GetFloatValueFromCache(key string) (retorno float64) {
	a := GetValueFromCache(key)
	if a != nil {
		retorno = a.(float64)
	}
	return
}

/*
GetStringValueFromCache gets an string value from the system's cache
*/
func GetStringValueFromCache(key string) (retorno string) {
	a := GetValueFromCache(key)
	if a != nil {
		retorno = a.(string)
	}
	return
}