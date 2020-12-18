package paramhandler

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/EricChiou/httprouter"
	"github.com/valyala/fasthttp"
)

// Set params from request ctx
func Set(context *httprouter.Context, varPtr interface{}, checkReq ...bool) error {
	check := true
	if len(checkReq) > 0 {
		check = checkReq[0]
	}

	// check variable is ptr
	if reflect.ValueOf(varPtr).Kind() != reflect.Ptr || varPtr == nil {
		return errors.New("need to input ptr")
	}

	// check variable type
	if reflect.ValueOf(varPtr).Elem().Kind() != reflect.Struct {
		return errors.New("need to input struct type")
	}

	// check path params
	for _, param := range context.Params {
		if len(param.Value) == 0 {
			return errors.New("path param format error")
		}
	}

	// parse body
	if body := context.Ctx.PostBody(); len(body) > 0 {
		if err := json.Unmarshal(body, varPtr); err != nil {
			return errors.New("body json parse error")
		}
	}

	// set parameters
	if err := setParam(context.Ctx, varPtr); err != nil {
		return err
	}

	// check required parameters
	if check {
		if err := checkParam(varPtr); err != nil {
			return err
		}
	}

	return nil
}

func setParam(ctx *fasthttp.RequestCtx, varPtr interface{}) error {
	reflectTyp := reflect.TypeOf(varPtr).Elem()
	reflectVal := reflect.ValueOf(varPtr).Elem()

	for i := 0; i < reflectVal.NumField(); i++ {
		jsonTag := reflectTyp.Field(i).Tag.Get("json")
		key := strings.Split(jsonTag, ",")[0]

		if len(key) > 0 {
			if varStr := string(ctx.QueryArgs().Peek(key)); len(varStr) > 0 {
				switch reflectVal.Field(i).Kind() {
				case reflect.Slice, reflect.Array:
					if err := setArray(varStr, reflectVal.Field(i)); err != nil {
						return err
					}

				case reflect.Struct:
					if reflectVal.Field(i).CanInterface() {
						variable := reflectVal.Field(i).Interface()
						if err := json.Unmarshal([]byte(varStr), &variable); err != nil {
							return errors.New("param format error")
						}

						if reflectVal.Field(i).CanSet() {
							reflectVal.Field(i).Set(reflect.ValueOf(&variable).Elem())

						} else {
							return errors.New("param format error")
						}

					} else {
						return errors.New("param format error")
					}

				default:
					if err := setVal(varStr, reflectVal.Field(i)); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func setVal(varStr string, reflectVal reflect.Value) error {
	switch reflectVal.Kind() {
	case reflect.String:
		reflectVal.SetString(varStr)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(varStr, 10, 64)
		if err != nil {
			return err
		}
		reflectVal.SetInt(num)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(varStr, 10, 64)
		if err != nil {
			return err
		}
		reflectVal.SetUint(num)

	case reflect.Float32, reflect.Float64:
		num, err := strconv.ParseFloat(varStr, 64)
		if err != nil {
			return err
		}
		reflectVal.SetFloat(num)

	case reflect.Complex64, reflect.Complex128:
		num, err := strconv.ParseFloat(varStr, 64)
		if err != nil {
			return err
		}
		reflectVal.SetComplex(complex(num, 0))

	case reflect.Bool:
		boolean, err := strconv.ParseBool(varStr)
		if err != nil {
			return err
		}
		reflectVal.SetBool(boolean)
	}

	return nil
}

func setArray(varStr string, reflectVal reflect.Value) error {
	varStrAry := strings.Split(varStr, ",")
	slice := reflect.MakeSlice(reflectVal.Type(), len(varStrAry), len(varStrAry))
	reflectVal.Set(slice)

	if reflectVal.Len() > 0 {
		switch reflectVal.Index(0).Kind() {
		case reflect.String:
			for i := 0; i < len(varStrAry); i++ {
				reflectVal.Index(i).SetString(varStrAry[i])
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for i := 0; i < len(varStrAry); i++ {
				num, err := strconv.ParseInt(varStrAry[i], 10, 64)
				if err != nil {
					return err
				}
				reflectVal.Index(i).SetInt(num)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for i := 0; i < len(varStrAry); i++ {
				num, err := strconv.ParseUint(varStrAry[i], 10, 64)
				if err != nil {
					return err
				}
				reflectVal.Index(i).SetUint(num)
			}

		case reflect.Float32, reflect.Float64:
			for i := 0; i < len(varStrAry); i++ {
				num, err := strconv.ParseFloat(varStrAry[i], 64)
				if err != nil {
					return err
				}
				reflectVal.Index(i).SetFloat(num)
			}

		case reflect.Complex64, reflect.Complex128:
			for i := 0; i < len(varStrAry); i++ {
				num, err := strconv.ParseFloat(varStrAry[i], 64)
				if err != nil {
					return err
				}
				reflectVal.Index(i).SetComplex(complex(num, 0))
			}

		case reflect.Bool:
			for i := 0; i < len(varStrAry); i++ {
				boolean, err := strconv.ParseBool(varStrAry[i])
				if err != nil {
					return err
				}
				reflectVal.Index(i).SetBool(boolean)
			}
		}
	}

	return nil
}

func checkParam(varPtr interface{}) error {
	reflectTyp := reflect.TypeOf(varPtr).Elem()
	reflectVal := reflect.ValueOf(varPtr).Elem()

	switch reflectVal.Kind() {
	case reflect.String:
		if reflectVal.CanInterface() {
			val := reflectVal.Interface()
			if len(val.(string)) == 0 {
				return errors.New("param format error")
			}

		} else {
			return errors.New("param format error")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		if reflectVal.CanInterface() {
			val := reflectVal.Interface()
			if val == 0 {
				return errors.New("param format error")
			}

		} else {
			return errors.New("param format error")
		}

	case reflect.Bool:
		return nil

	case reflect.Slice, reflect.Array:
		for i := 0; i < reflectVal.Len(); i++ {
			if reflectVal.Index(i).Kind() == reflect.Struct {
				if reflectVal.Index(i).CanAddr() && reflectVal.Index(i).Addr().CanInterface() {
					if err := checkParam(reflectVal.Index(i).Addr().Interface()); err != nil {
						return err
					}

				} else {
					return errors.New("param format error")
				}
			}
		}

	case reflect.Struct:
		for i := 0; i < reflectVal.NumField(); i++ {
			req, _ := strconv.ParseBool(reflectTyp.Field(i).Tag.Get("req"))
			if req {
				if reflectVal.Field(i).CanAddr() && reflectVal.Field(i).Addr().CanInterface() {
					if err := checkParam(reflectVal.Field(i).Addr().Interface()); err != nil {
						return err
					}

				} else {
					return errors.New("param format error")
				}
			}
		}
	}

	return nil
}
