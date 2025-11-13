package guia2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Element struct {
	parent *Driver
	id     string
}

func (e *Element) Text() (text string, err error) {
	// register(getHandler, new GetText("/wd/hub/session/:sessionId/element/:id/text"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/text"); err != nil {
		return "", err
	}
	var reply = new(struct{ Value string })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return "", err
	}
	text = reply.Value
	return
}

func (e *Element) GetAttribute(name string) (attribute string, err error) {
	// register(getHandler, new GetElementAttribute("/wd/hub/session/:sessionId/element/:id/attribute/:name"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/attribute", name); err != nil {
		return "", err
	}
	var reply = new(struct{ Value string })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return "", err
	}
	attribute = reply.Value
	return
}

func (e *Element) ContentDescription() (name string, err error) {
	// register(getHandler, new GetName("/wd/hub/session/:sessionId/element/:id/name"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/name"); err != nil {
		return "", err
	}
	var reply = new(struct{ Value string })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return "", err
	}
	name = reply.Value
	return
}

func (e *Element) Size() (size Size, err error) {
	// register(getHandler, new GetSize("/wd/hub/session/:sessionId/element/:id/size"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/size"); err != nil {
		return Size{-1, -1}, err
	}
	var reply = new(struct{ Value Size })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Size{-1, -1}, err
	}
	size = reply.Value
	return
}

type Rect struct {
	Point
	Size
}

func (e *Element) Rect() (rect Rect, err error) {
	// register(getHandler, new GetRect("/wd/hub/session/:sessionId/element/:id/rect"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/rect"); err != nil {
		return Rect{}, err
	}
	var reply = new(struct{ Value Rect })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Rect{}, err
	}
	rect = reply.Value
	return
}

func (e *Element) Screenshot() (raw *bytes.Buffer, err error) {
	// W3C endpoint
	// register(getHandler, new GetElementScreenshot("/wd/hub/session/:sessionId/element/:id/screenshot"))
	// JSONWP endpoint
	// register(getHandler, new GetElementScreenshot("/wd/hub/session/:sessionId/screenshot/:id"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/screenshot"); err != nil {
		return nil, err
	}
	var reply = new(struct{ Value string })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, err
	}

	var decodeStr []byte
	if decodeStr, err = base64.StdEncoding.DecodeString(reply.Value); err != nil {
		return nil, err
	}

	raw = bytes.NewBuffer(decodeStr)
	return
}

func (e *Element) Location() (point Point, err error) {
	// register(getHandler, new Location("/wd/hub/session/:sessionId/element/:id/location"))
	var rawResp RawResponse
	if rawResp, err = e.parent.executeGet("/session", e.parent.sessionId, "/element", e.id, "/location"); err != nil {
		return Point{-1, -1}, err
	}
	var reply = new(struct{ Value Point })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Point{-1, -1}, err
	}
	point = reply.Value
	return
}

func (e *Element) Click() (err error) {
	// register(postHandler, new Click("/wd/hub/session/:sessionId/element/:id/click"))
	_, err = e.parent.executePost(nil, "/session", e.parent.sessionId, "/element", e.id, "/click")
	return
}

func (e *Element) DoubleClick() error {
	payload := map[string]interface{}{
		"origin": makeElementRef(e.id),
	}
	return e.parent.performGesture("double_click", payload)
}

func (e *Element) LongClick(duration ...float64) error {
	payload := map[string]interface{}{
		"origin": makeElementRef(e.id),
	}
	if len(duration) != 0 {
		payload["duration"] = int64(duration[0] * 1000)
	}
	return e.parent.performGesture("long_click", payload)
}

func (e *Element) Clear() (err error) {
	// register(postHandler, new Clear("/wd/hub/session/:sessionId/element/:id/clear"))
	_, err = e.parent.executePost(nil, "/session", e.parent.sessionId, "/element", e.id, "/clear")
	return
}

func (e *Element) SendKeys(text string, isReplace ...bool) (err error) {
	if len(isReplace) == 0 {
		isReplace = []bool{true}
	}
	// register(postHandler, new SendKeysToElement("/wd/hub/session/:sessionId/element/:id/value"))
	// https://github.com/appium/appium-uiautomator2-server/blob/master/app/src/main/java/io/appium/uiautomator2/handler/SendKeysToElement.java#L76-L85
	data := map[string]interface{}{
		"text":    text,
		"replace": isReplace[0],
	}
	_, err = e.parent.executePost(data, "/session", e.parent.sessionId, "/element", e.id, "/value")
	return
}

func (e *Element) FindElements(by BySelector) (elements []*Element, err error) {
	method, selector := by.getMethodAndSelector()
	return e.parent._findElements(method, selector, e.id)
}

func (e *Element) FindElement(by BySelector) (elem *Element, err error) {
	method, selector := by.getMethodAndSelector()
	return e.parent._findElement(method, selector, e.id)
}

func (e *Element) Swipe(startX, startY, endX, endY int, steps ...int) (err error) {
	return e.SwipeFloat(float64(startX), float64(startY), float64(endX), float64(endY), steps...)
}

func (e *Element) SwipeFloat(startX, startY, endX, endY float64, steps ...int) (err error) {
	rect, err := e.Rect()
	if err != nil {
		return err
	}
	start := PointF{
		X: float64(rect.Point.X) + startX,
		Y: float64(rect.Point.Y) + startY,
	}
	end := PointF{
		X: float64(rect.Point.X) + endX,
		Y: float64(rect.Point.Y) + endY,
	}
	return e.parent.DragFloat(start.X, start.Y, end.X, end.Y, steps...)
}

func (e *Element) SwipePoint(startPoint, endPoint Point, steps ...int) (err error) {
	return e.Swipe(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y, steps...)
}

func (e *Element) SwipePointF(startPoint, endPoint PointF, steps ...int) (err error) {
	return e.SwipeFloat(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y, steps...)
}

func (e *Element) Drag(endX, endY int, steps ...int) (err error) {
	return e.DragFloat(float64(endX), float64(endY), steps...)
}

func (e *Element) DragFloat(endX, endY float64, steps ...int) error {
	req := dragRequest{
		Origin: makeElementRef(e.id),
		End:    PointF{X: endX, Y: endY},
	}
	if len(steps) != 0 {
		speed := steps[0]
		req.Speed = &speed
	}
	return e.parent.sendDrag(req)
}

func (e *Element) DragPoint(endPoint Point, steps ...int) error {
	return e.Drag(endPoint.X, endPoint.Y, steps...)
}

func (e *Element) DragPointF(endPoint PointF, steps ...int) (err error) {
	return e.DragFloat(endPoint.X, endPoint.Y, steps...)
}

func (e *Element) DragTo(destElem *Element, steps ...int) error {
	rect, err := destElem.Rect()
	if err != nil {
		return err
	}
	centerX := float64(rect.Point.X + rect.Size.Width/2)
	centerY := float64(rect.Point.Y + rect.Size.Height/2)
	return e.DragFloat(centerX, centerY, steps...)
}

func (e *Element) Flick(xOffset, yOffset, speed int) (err error) {
	if xOffset == 0 && yOffset == 0 {
		return errors.New("both 'xOffset' and 'yOffset' cannot be zero")
	}
	direction, err := directionFromVector(xOffset, yOffset)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"origin":    makeElementRef(e.id),
		"direction": string(direction),
	}
	if speed != 0 {
		if speed < 0 {
			speed = -speed
		}
		payload["speed"] = speed
	}
	var completed bool
	if err := e.parent.postGestureForValue("fling", payload, &completed); err != nil {
		return err
	}
	if !completed {
		return errors.New("fling gesture did not complete")
	}
	return nil
}

func (e *Element) ScrollTo(by BySelector, maxSwipes ...int) (err error) {
	if len(maxSwipes) == 0 {
		maxSwipes = []int{0}
	}
	method, selector := by.getMethodAndSelector()
	return e.parent._scrollTo(method, selector, maxSwipes[0], e.id)
}

func (e *Element) ScrollToElement(element *Element) (err error) {
	// register(postHandler, new ScrollToElement("/wd/hub/session/:sessionId/appium/element/:id/scroll_to/:id2"))
	_, err = e.parent.executePost(nil, "/session", e.parent.sessionId, "/appium/element", e.id, "/scroll_to", element.id)
	return
}
