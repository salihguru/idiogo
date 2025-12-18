package xopt

func Get[T any](def T, opts ...T) T {
	if len(opts) == 0 {
		return def
	}
	return opts[0]
}

func GetNotEmpty[T any](opt *T, def T) T {
	if opt == nil {
		return def
	}
	return *opt
}

func GetNotEmptyStr(opt string, def string) string {
	if opt == "" {
		return def
	}
	return opt
}
