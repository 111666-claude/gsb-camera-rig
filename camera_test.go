package camera

import "testing"

func TestNewStartsAtWorldCenter(t *testing.T) {
	cam := New(Config{WorldW: 1000, WorldH: 800, ViewW: 400, ViewH: 300, Zoom: 1})
	x, y := cam.Center()
	if x != 500 || y != 400 {
		t.Fatalf("初始应该在 (500,400)：%v %v", x, y)
	}
}

func TestViewUsesZoom(t *testing.T) {
	cam := New(Config{WorldW: 1000, WorldH: 1000, ViewW: 800, ViewH: 600, Zoom: 2})
	view := cam.View()
	if view.W != 400 || view.H != 300 {
		t.Fatalf("缩放后视野应该是 400x300：%v %v", view.W, view.H)
	}
}

func TestTeleportMovesCenter(t *testing.T) {
	cam := New(Config{WorldW: 1000, WorldH: 1000, ViewW: 100, ViewH: 100, Zoom: 1})
	cam.Teleport(120, 340)
	x, y := cam.Center()
	if x != 120 || y != 340 {
		t.Fatalf("Teleport 之后中心应该是 (120,340)：%v %v", x, y)
	}
}

func TestWorkStartsAtZero(t *testing.T) {
	if New(Config{WorldW: 10, WorldH: 10, ViewW: 10, ViewH: 10, Zoom: 1}).Work() != 0 {
		t.Fatal("Work 初始应该是 0")
	}
}

func TestSmoothFactorBounds(t *testing.T) {
	if got := SmoothFactor(0, 10); got != 1 {
		t.Fatalf("没有平滑时间时系数应该是 1：%v", got)
	}
	if got := SmoothFactor(200, 0); got != 0 {
		t.Fatalf("零时长系数应该是 0：%v", got)
	}
}
