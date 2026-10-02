// Command camera-rig 跑相机样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/camera"
)

func baseConfig() camera.Config {
	return camera.Config{
		WorldW: 4000, WorldH: 3000,
		ViewW: 800, ViewH: 600,
		Zoom: 1, DeadW: 100, DeadH: 100, SmoothMs: 0,
	}
}

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("camera-rig", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "deadzone", "deadzone / clamp / fps / smallworld / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "deadzone":
		cam := camera.New(baseConfig())
		cam.Teleport(500, 500)
		cam.Update(540, 520, 16)
		x, y := cam.Center()
		fmt.Fprintf(stdout, "center=%.0f,%.0f\n", x, y)
	case "clamp":
		cfg := baseConfig()
		cfg.Zoom = 2
		cam := camera.New(cfg)
		cam.Teleport(500, 500)
		cam.Update(3900, 2900, 16)
		x, y := cam.Center()
		fmt.Fprintf(stdout, "center=%.0f,%.0f\n", x, y)
	case "fps":
		cfg := baseConfig()
		cfg.SmoothMs = 200
		fast := camera.New(cfg)
		fast.Teleport(0, 0)
		for index := 0; index < 30; index++ {
			fast.Update(1000, 0, 4)
		}
		slow := camera.New(cfg)
		slow.Teleport(0, 0)
		for index := 0; index < 8; index++ {
			slow.Update(1000, 0, 15)
		}
		ax, _ := fast.Center()
		bx, _ := slow.Center()
		fmt.Fprintf(stdout, "a=%.1f b=%.1f diff=%.1f\n", ax, bx, ax-bx)
	case "smallworld":
		cfg := baseConfig()
		cfg.WorldW, cfg.WorldH = 400, 300
		cam := camera.New(cfg)
		cam.Teleport(200, 150)
		cam.Update(50, 50, 16)
		x, y := cam.Center()
		fmt.Fprintf(stdout, "center=%.0f,%.0f\n", x, y)
	case "scan":
		cam := camera.New(baseConfig())
		for index := 0; index < 3000; index++ {
			cam.AddChunk(camera.Rect{X: float64(index), W: 10, H: 10})
		}
		for index := 0; index < 3000; index++ {
			cam.Update(1000, 1000, 16)
		}
		fmt.Fprintf(stdout, "work=%d\n", cam.Work())
	default:
		fmt.Fprintln(stderr, "需要 --sample deadzone|clamp|fps|smallworld|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
