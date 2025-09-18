package core

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

var (
	// TimeToStringConverter is a list of type converters for time to string.
	TimeToStringConverter = []copier.TypeConverter{
		{
			SrcType: (*time.Time)(nil),
			DstType: (*string)(nil),
			Fn: func(src any) (any, error) {
				t, ok := src.(*time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if t != nil {
					s := t.Format(time.DateTime)
					return &s, nil
				}
				return (*string)(nil), nil
			},
		},
		{
			SrcType: time.Time{},
			DstType: "",
			Fn: func(src any) (any, error) {
				t, ok := src.(time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				return t.Format(time.DateTime), nil
			},
		},
	}

	// StringToTimeConverter is a list of type converters for string to time.
	StringToTimeConverter = []copier.TypeConverter{
		{
			SrcType: (*string)(nil),
			DstType: (*time.Time)(nil),
			Fn: func(src any) (any, error) {
				s, ok := src.(*string)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s != nil && *s != "" {
					t, err := time.Parse(time.DateTime, *s)
					if err != nil {
						return nil, err
					}
					return &t, nil
				}
				return (*time.Time)(nil), nil
			},
		},
		{
			SrcType: "",
			DstType: time.Time{},
			Fn: func(src any) (any, error) {
				s, ok := src.(string)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s != "" {
					return time.Parse(time.DateTime, s)
				}
				return time.Time{}, nil
			},
		},
	}
)

// TypeConverters returns a list of type converters for copier.
func TypeConverters() []copier.TypeConverter {
	var converters []copier.TypeConverter

	converters = append(converters, TimeToStringConverter...)
	converters = append(converters, StringToTimeConverter...)
	return converters
}

// CopyWithConverters copies the value from to to the value from from with the type converters.
func CopyWithConverters(to any, from any) error {
	return copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters:  TypeConverters(),
	})
}

// Copy copies the value from to to the value from from with the default type converters.
func Copy(to any, from any) error {
	return copier.Copy(to, from)
}
