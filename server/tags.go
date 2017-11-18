package server

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/flosch/pongo2"
)

// lookup retrieves the value for the specified field if the item is a struct
// or the value of the specified key if the item is a map.
func lookup(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var (
		vVal = reflect.ValueOf(param.Interface())
		key  = in.String()
		v    reflect.Value
	)
	for vVal.Kind() == reflect.Ptr {
		vVal = vVal.Elem()
	}
	switch vVal.Kind() {
	case reflect.Struct:
		v = vVal.FieldByName(key)
	case reflect.Map:
		v = vVal.MapIndex(reflect.ValueOf(key))
		if !v.IsValid() {
			return pongo2.AsValue(nil), nil
		}
	default:
		return nil, &pongo2.Error{
			Sender:    "filter:lookup",
			OrigError: errors.New("parameter is not a struct or map"),
		}
	}
	return pongo2.AsValue(v.Interface()), nil
}

// tagNode stores the information needed to evaluate the tag tag.
type tagNode struct {
	position  *pongo2.Token
	varName   string
	structVal pongo2.IEvaluator
	fieldName pongo2.IEvaluator
	tagName   pongo2.IEvaluator
}

func (t *tagNode) Execute(ctx *pongo2.ExecutionContext, w pongo2.TemplateWriter) *pongo2.Error {
	structVal, err := t.structVal.Evaluate(ctx)
	if err != nil {
		return err
	}
	fieldName, err := t.fieldName.Evaluate(ctx)
	if err != nil {
		return err
	}
	tagName, err := t.tagName.Evaluate(ctx)
	if err != nil {
		return err
	}
	fType, ok := reflect.TypeOf(structVal.Interface()).FieldByName(fieldName.String())
	if !ok {
		return ctx.Error(fmt.Sprintf("\"%s\" is not a field", fieldName), t.position)
	}
	ctx.Private[t.varName] = pongo2.AsValue(fType.Tag.Get(tagName.String()))
	return nil
}

// tag retrieves the value for the specified tag in the specified field.
func tag(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	nameToken := arguments.MatchType(pongo2.TokenIdentifier)
	if nameToken == nil {
		return nil, arguments.Error("identifier expected", start)
	}
	if arguments.Match(pongo2.TokenSymbol, "=") == nil {
		return nil, arguments.Error("\"=\" expected", start)
	}
	structVal, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	fieldName, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	tagName, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	if arguments.Remaining() > 0 {
		return nil, arguments.Error("unexpected arguments", start)
	}
	return &tagNode{
		position:  start,
		varName:   nameToken.Val,
		structVal: structVal,
		fieldName: fieldName,
		tagName:   tagName,
	}, nil
}

func init() {
	pongo2.RegisterFilter("lookup", lookup)
	pongo2.RegisterTag("tag", tag)
}
