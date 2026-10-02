// Package camera 是跟随相机：死区、边界钳制、缩放换算与按真实时间平滑。
// 缺陷：没有死区、钳制用未缩放的视野、插值按帧而不是按时间、世界比视野小时不居中、
// 每次更新都扫描全部分块。
package camera

import "math"

// Config 是相机配置。
type Config struct {
	WorldW   float64
	WorldH   float64
	ViewW    float64
	ViewH    float64
	Zoom     float64
	DeadW    float64
	DeadH    float64
	SmoothMs float64
}

// Rect 是一块世界区域。
type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

// Camera 是跟随相机。
type Camera struct {
	cfg      Config
	cx       float64
	cy       float64
	chunks   []Rect
	work     int
	lastGoal [2]float64
	hasGoal  bool
	anchorX  float64
	anchorY  float64
	anchorMs float64
	nowMs    float64
}

// New 建相机，初始停在世界中心。
func New(cfg Config) *Camera {
	if cfg.Zoom == 0 {
		cfg.Zoom = 1
	}
	return &Camera{cfg: cfg, cx: cfg.WorldW / 2, cy: cfg.WorldH / 2}
}

// AddChunk 登记一块世界分块。
func (c *Camera) AddChunk(r Rect) { c.chunks = append(c.chunks, r) }

// Teleport 直接把相机放到某个位置。
func (c *Camera) Teleport(x, y float64) { c.cx, c.cy = x, y }

// Update 推进一帧。
// 缺陷：每次更新扫描全部分块；插值系数只跟单帧时长成正比，和帧率相关；
// 平滑没有从锚点按累计时间重算。
func (c *Camera) Update(tx, ty, dtMs float64) {
	c.work += len(c.chunks)
	gx, gy := c.goal(tx, ty)
	factor := dtMs / c.cfg.SmoothMs
	if c.cfg.SmoothMs <= 0 {
		factor = 1
	}
	c.cx += (gx - c.cx) * factor
	c.cy += (gy - c.cy) * factor
}

// goal 算这一帧的目标中心。
// 缺陷：没有死区；钳制用未按缩放换算的视野。
func (c *Camera) goal(tx, ty float64) (float64, float64) {
	gx, gy := tx, ty
	halfW := c.cfg.ViewW / 2
	halfH := c.cfg.ViewH / 2
	if gx < halfW {
		gx = halfW
	}
	if gx > c.cfg.WorldW-halfW {
		gx = c.cfg.WorldW - halfW
	}
	if gy < halfH {
		gy = halfH
	}
	if gy > c.cfg.WorldH-halfH {
		gy = c.cfg.WorldH - halfH
	}
	return gx, gy
}

// Center 是当前中心。
func (c *Camera) Center() (float64, float64) { return c.cx, c.cy }

// View 是当前可见的世界矩形。
func (c *Camera) View() Rect {
	halfW := c.cfg.ViewW / c.cfg.Zoom / 2
	halfH := c.cfg.ViewH / c.cfg.Zoom / 2
	return Rect{X: c.cx - halfW, Y: c.cy - halfH, W: halfW * 2, H: halfH * 2}
}

// Work 是累计扫描量（规模观测）。
func (c *Camera) Work() int { return c.work }

// SmoothFactor 是给定累计时长后的平滑系数（参考用）。
func SmoothFactor(smoothMs, elapsedMs float64) float64 {
	if smoothMs <= 0 {
		return 1
	}
	return 1 - math.Exp(-elapsedMs/smoothMs)
}
