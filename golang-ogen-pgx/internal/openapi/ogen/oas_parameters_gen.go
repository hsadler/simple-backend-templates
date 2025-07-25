// Code generated by ogen, DO NOT EDIT.

package ogen

import (
	"net/http"
	"net/url"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// ItemsAllGetParams is parameters of GET /items/all operation.
type ItemsAllGetParams struct {
	// Offset.
	Offset int
	// Chunk size.
	ChunkSize int
}

func unpackItemsAllGetParams(packed middleware.Parameters) (params ItemsAllGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "offset",
			In:   "query",
		}
		params.Offset = packed[key].(int)
	}
	{
		key := middleware.ParameterKey{
			Name: "chunkSize",
			In:   "query",
		}
		params.ChunkSize = packed[key].(int)
	}
	return params
}

func decodeItemsAllGetParams(args [0]string, argsEscaped bool, r *http.Request) (params ItemsAllGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: offset.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "offset",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.Offset = c
				return nil
			}); err != nil {
				return err
			}
			if err := func() error {
				if err := (validate.Int{
					MinSet:        true,
					Min:           0,
					MaxSet:        false,
					Max:           0,
					MinExclusive:  false,
					MaxExclusive:  false,
					MultipleOfSet: false,
					MultipleOf:    0,
				}).Validate(int64(params.Offset)); err != nil {
					return errors.Wrap(err, "int")
				}
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "offset",
			In:   "query",
			Err:  err,
		}
	}
	// Decode query: chunkSize.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "chunkSize",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.ChunkSize = c
				return nil
			}); err != nil {
				return err
			}
			if err := func() error {
				if err := (validate.Int{
					MinSet:        true,
					Min:           1,
					MaxSet:        true,
					Max:           20,
					MinExclusive:  false,
					MaxExclusive:  false,
					MultipleOfSet: false,
					MultipleOf:    0,
				}).Validate(int64(params.ChunkSize)); err != nil {
					return errors.Wrap(err, "int")
				}
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "chunkSize",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// ItemsGetParams is parameters of GET /items operation.
type ItemsGetParams struct {
	// Item IDs.
	ItemIds []int
}

func unpackItemsGetParams(packed middleware.Parameters) (params ItemsGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "item_ids",
			In:   "query",
		}
		params.ItemIds = packed[key].([]int)
	}
	return params
}

func decodeItemsGetParams(args [0]string, argsEscaped bool, r *http.Request) (params ItemsGetParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: item_ids.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "item_ids",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				return d.DecodeArray(func(d uri.Decoder) error {
					var paramsDotItemIdsVal int
					if err := func() error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToInt(val)
						if err != nil {
							return err
						}

						paramsDotItemIdsVal = c
						return nil
					}(); err != nil {
						return err
					}
					params.ItemIds = append(params.ItemIds, paramsDotItemIdsVal)
					return nil
				})
			}); err != nil {
				return err
			}
			if err := func() error {
				if params.ItemIds == nil {
					return errors.New("nil is invalid value")
				}
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "item_ids",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// ItemsIDGetParams is parameters of GET /items/{id} operation.
type ItemsIDGetParams struct {
	// Item ID.
	ID int
}

func unpackItemsIDGetParams(packed middleware.Parameters) (params ItemsIDGetParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "path",
		}
		params.ID = packed[key].(int)
	}
	return params
}

func decodeItemsIDGetParams(args [1]string, argsEscaped bool, r *http.Request) (params ItemsIDGetParams, _ error) {
	// Decode path: id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.ID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
