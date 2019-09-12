package dt

import (
	"errors"
	"strconv"
	"time"
)

// GetInt чтение поля записи по имени
func (r *Data) GetInt(name string) (int, error) {
	val, ok := r.values[name]
	if ok {
		return strconv.Atoi(val)
	}
	err := errors.New("Not found field " + name)
	return 0, err
}

// GetBool чтение поля записи по имени
func (r *Data) GetBool(name string) (bool, error) {
	val, ok := r.values[name]
	if ok {
		return strconv.ParseBool(val)
	}
	return false, errors.New("Not found field " + name)
}

// GetFloat чтение поля записи по имени
func (r *Data) GetFloat(name string) (float64, error) {
	val, ok := r.values[name]
	if ok {
		return strconv.ParseFloat(val, 64)
	}
	return 0.0, errors.New("Not found field " + name)
}

// GetLong чтение поля записи по имени
func (r *Data) GetLong(name string) (int64, error) {
	val, ok := r.values[name]
	if ok {
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, errors.New("Not found field " + name)
}

// GetString чтение поля записи по имени
func (r *Data) GetString(name string) (ret string, err error) {
	val, ok := r.values[name]
	if ok {
		return val, nil
	}
	err = errors.New("Not found field " + name)
	return ret, err

}

// GetDate чтение поля записи по имени
func (r *Data) GetDate(name string) (time.Time, error) {
	val, ok := r.values[name]
	if ok {
		return time.Parse(time.RFC3339Nano, val)
	}
	return time.Now(), errors.New("Not found field " + name)
}

//SetInt функция
func (r *Data) SetInt(name string, value int) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = strconv.Itoa(value)
		return nil
	}
	return errors.New("Not found field " + name)
}

//SetBool функция
func (r *Data) SetBool(name string, value bool) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = strconv.FormatBool(value)
		return nil
	}
	return errors.New("Not found field " + name)
}

//SetFloat функция
func (r *Data) SetFloat(name string, value float64) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = strconv.FormatFloat(value, 'f', 5, 64)
		return nil
	}
	return errors.New("Not found field " + name)
}

//SetLong функция
func (r *Data) SetLong(name string, value int64) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = strconv.FormatInt(value, 10)
		return nil
	}
	return errors.New("Not found field " + name)
}

//SetString функция
func (r *Data) SetString(name string, value string) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = value
		return nil
	}
	return errors.New("Not found  field " + name)
}

//SetDate функция
func (r *Data) SetDate(name string, value time.Time) error {
	_, ok := r.values[name]
	if ok {
		r.values[name] = value.Format(time.RFC3339Nano)
		return nil
	}
	return errors.New("Not found")

}
