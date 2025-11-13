package guia2

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

// RectF represents a floating-point rectangle that can be sent to the new Appium gesture endpoints.
type RectF struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func newRectF(left, top, width, height float64) RectF {
	return RectF{
		Left:   left,
		Top:    top,
		Width:  width,
		Height: height,
	}
}

func rectFFromSize(size Size) RectF {
	return newRectF(0, 0, float64(size.Width), float64(size.Height))
}

func rectFFromRect(rect Rect) RectF {
	return newRectF(
		float64(rect.Point.X),
		float64(rect.Point.Y),
		float64(rect.Size.Width),
		float64(rect.Size.Height),
	)
}

// GestureDirection is the direction string expected by the new gesture endpoints.
type GestureDirection string

const (
	DirectionUp    GestureDirection = "up"
	DirectionDown  GestureDirection = "down"
	DirectionLeft  GestureDirection = "left"
	DirectionRight GestureDirection = "right"
)

var errLegacyTouchCommand = errors.New("legacy JSONWP touch actions have been removed in appium-uiautomator2-server >= 7.0; please migrate to W3C actions")

func makeElementRef(id string) map[string]string {
	if id == "" {
		return nil
	}
	return map[string]string{
		legacyWebElementIdentifier: id,
		webElementIdentifier:       id,
	}
}

func (d *Driver) gesturePath(action string) string {
	return fmt.Sprintf("appium/gestures/%s", action)
}

func (d *Driver) performGesture(action string, payload interface{}) error {
	return d.postGestureForValue(action, payload, nil)
}

func (d *Driver) postGestureForValue(action string, payload interface{}, out interface{}) error {
	rawResp, err := d.executePost(payload, "/session", d.sessionId, d.gesturePath(action))
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	return decodeValue(rawResp, out)
}

func decodeValue(raw RawResponse, out interface{}) error {
	if out == nil {
		return nil
	}
	var envelope struct {
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return err
	}
	if len(envelope.Value) == 0 {
		return nil
	}
	return json.Unmarshal(envelope.Value, out)
}

func directionFromVector(xDelta, yDelta int) (GestureDirection, error) {
	if xDelta == 0 && yDelta == 0 {
		return "", errors.New("direction cannot be determined from a zero vector")
	}
	if math.Abs(float64(xDelta)) >= math.Abs(float64(yDelta)) {
		if xDelta > 0 {
			return DirectionRight, nil
		}
		return DirectionLeft, nil
	}
	if yDelta > 0 {
		return DirectionDown, nil
	}
	return DirectionUp, nil
}

func (d *Driver) fullScreenArea() (RectF, error) {
	size, err := d.DeviceSize()
	if err != nil {
		return RectF{}, err
	}
	return rectFFromSize(size), nil
}
