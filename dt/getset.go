package dt

import "errors"
import "time"

// GetInt чтение поля записи по имени
func (r *Data) GetInt(name string) (ret int, err error) {
	val, ok := r.values[name]
	if ok {
		return val.getInt(), nil
	}
	err = errors.New("Not found")
	return ret, err
}

// GetValue чтение поля записи по имени
func (r *Data) GetValue(name string) (ret string, err error) {
	val, ok := r.values[name]
	if ok {
		return val.value, nil
	}
	err = errors.New("Not found")
	return ret, err
}

// GetBool чтение поля записи по имени
func (r *Data) GetBool(name string) (bool, error) {
	val, ok := r.values[name]
	if ok {
		return val.getBool(), nil
	}
	return false, errors.New("Not found")
}

// GetFloat чтение поля записи по имени
func (r *Data) GetFloat(name string) (float64, error) {
	val, ok := r.values[name]
	if ok {
		return val.getFloat(), nil
	}
	return 0.0, errors.New("Not found")
}

// GetLong чтение поля записи по имени
func (r *Data) GetLong(name string) (int64, error) {
	val, ok := r.values[name]
	if ok {
		return val.getLong(), nil
	}
	return 0, errors.New("Not found")
}

// GetString чтение поля записи по имени
func (r *Data) GetString(name string) (ret string, err error) {
	val, ok := r.values[name]
	if ok {
		return val.getString(), nil
	}
	err = errors.New("Not found")
	return ret, err

}

// GetDate чтение поля записи по имени
func (r *Data) GetDate(name string) (time.Time, error) {
	val, ok := r.values[name]
	if ok {
		return val.getDate(), nil
	}
	return time.Now(), errors.New("Not found")
}

//SetInt функция
func (r *Data) SetInt(name string, value int) error {
	val, ok := r.values[name]
	if ok {
		val.setInt(value)
		return nil
	}
	return errors.New("Not found")
}

//SetBool функция
func (r *Data) SetBool(name string, value bool) error {
	val, ok := r.values[name]
	if ok {
		val.setBool(value)
		return nil
	}
	return errors.New("Not found")
}

//SetFloat функция
func (r *Data) SetFloat(name string, value float64) error {
	val, ok := r.values[name]
	if ok {
		val.setFloat(value)
		return nil
	}
	return errors.New("Not found")
}

//SetLong функция
func (r *Data) SetLong(name string, value int64) error {
	val, ok := r.values[name]
	if ok {
		val.setLong(value)
		return nil
	}
	return errors.New("Not found")
}

//SetString функция
func (r *Data) SetString(name string, value string) error {
	val, ok := r.values[name]
	if ok {
		val.setString(value)
		return nil
	}
	return errors.New("Not found")
}

//SetDate функция
func (r *Data) SetDate(name string, value time.Time) error {
	val, ok := r.values[name]
	if ok {
		val.setDate(value)
		return nil
	}
	return errors.New("Not found")

}
