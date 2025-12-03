package rik

type JSONBuilder struct {
    data map[string]any
}

func NewJSON() *JSONBuilder {
    return &JSONBuilder{data: map[string]any{}}
}

func (j *JSONBuilder) Set(key string, value ...any) *JSONBuilder {
    valueLen := len(value)
    if valueLen == 0 {
        return j
    }
    if valueLen == 1 {
        j.data[key] = value[0]
    } else {
        j.data[key] = value
    }
    return j
}

func (j *JSONBuilder) Add(key string, value ...any) *JSONBuilder {
    if len(value) == 0 {
        return j
    }
    jsonArray := j.data[key]
    if jsonArray == nil {
        j.Set(key, value)
    } else {
        s, ok := jsonArray.([]any)
        if !ok {
            panic("cannot add value to this key. it is not slice")
        }
        j.data[key] = append(s, value...)
    }
    return j
}

func (j *JSONBuilder) SetAll(data map[string]any) *JSONBuilder {
    for k, v := range data {
        j.Set(k, v)
    }
    return j
}

func (j *JSONBuilder) AddAll(data map[string]any) *JSONBuilder {
    for k, v := range data {
        j.Add(k, v)
    }
    return j
}

func (j *JSONBuilder) Build() map[string]any {
    return j.data
}

func (j *JSONBuilder) BuildAsSingleArray() []map[string]any {
    return []map[string]any{j.data}
}
