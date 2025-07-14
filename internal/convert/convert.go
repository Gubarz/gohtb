package convert

func Slice[In any, Out any](input []In, fn func(In) Out) []Out {
	out := make([]Out, 0, len(input))
	for _, v := range input {
		out = append(out, fn(v))
	}
	return out
}

func Map[K comparable, V1 any, V2 any](m *map[K]V1, f func(V1) V2) map[K]V2 {
	if m == nil {
		return nil
	}
	result := make(map[K]V2, len(*m))
	for k, v := range *m {
		result[k] = f(v)
	}
	return result
}

func SliceCast[From, To any](in []From, convertFn func(From) To) []To {
	out := make([]To, len(in))
	for i, v := range in {
		out[i] = convertFn(v)
	}
	return out
}

func SlicePointer[In any, Out any](input *[]In, fn func(In) Out) []Out {
	if input == nil {
		return nil
	}
	return Slice(*input, fn)
}

func MapSlice[K comparable, V1 any, V2 any](m map[K]V1, f func(V1) V2) []V2 {
	if m == nil {
		return nil
	}
	out := make([]V2, 0, len(m))
	for _, v := range m {
		out = append(out, f(v))
	}
	return out
}
