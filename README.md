# camera-rig

跟随相机：死区、边界钳制、缩放换算与按真实时间平滑，只用 Go 标准库。

```
go run ./cmd/camera-rig --sample=deadzone
go run ./cmd/camera-rig --sample=clamp
go run ./cmd/camera-rig --sample=fps
go run ./cmd/camera-rig --sample=smallworld
go run ./cmd/camera-rig --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **缩放**：可见世界矩形的宽高等于 `ViewW/Zoom` 与 `ViewH/Zoom`。
- **死区**：目标与相机中心在某个轴上的偏差不超过死区的一半时，该轴的中心不动；
  超出时相机中心朝目标移动，正好让目标落在死区边缘上。
- **边界钳制**：可见矩形必须完全落在世界内；世界某个方向比可见矩形还小时，
  该轴中心固定在世界中点。
- **平滑**：中心从锚点按 `1-exp(-累计时长/SmoothMs)` 逼近目标；相同累计时长、
  不同帧步长必须得到相同结果；`SmoothMs<=0` 时直接到达目标。
- **代价**：`Update` 每帧 O(1)，`Work` 不随分块数量线性增长。
- **规模**：10 万块世界分块、每秒 10 万次更新，单次更新摊还 O(1)。

## 输出契约（不改格式）

```
center=..,..
a=.. b=.. diff=..
work=..
```

## 目录

```
camera.go              死区、钳制、缩放与平滑
camera_test.go         单元测试
cmd/camera-rig/        命令行入口
```
