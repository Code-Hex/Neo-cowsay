// +build windows

package super

func height() (int, error) {
	t, err := tty.Open()
	if err != nil {
		return 0, err
	}
	_, h, err := t.Size()
	if err != nil {
		return 0, err
	}
	return h, nil
}
