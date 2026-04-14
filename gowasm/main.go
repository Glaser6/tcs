//go:build js && wasm

package main

import (
	"math/rand"
	"strconv"
	"syscall/js"
	"time"
)

// 全局常量定义
const (
	canvasWidth  = 600 // 画布宽度
	canvasHeight = 600 // 画布高度
	gridSize     = 20  // 网格大小（蛇身/食物尺寸）
	updateSpeed  = 150 // 蛇移动间隔（毫秒）
)

// Direction 方向枚举
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Point 坐标结构体
type Point struct {
	X, Y int
}

// SnakeGame 游戏核心结构体
type SnakeGame struct {
	ctx       js.Value  // Canvas 2D上下文
	snake     []Point   // 蛇的身体
	direction Direction // 当前移动方向
	nextDir   Direction // 下一个移动方向（防止反向冲突）
	food      Point     // 食物坐标
	score     int       // 分数
	gameOver  bool      // 游戏结束标记
	canvas    js.Value  // Canvas元素
}

// NewSnakeGame 初始化游戏
func NewSnakeGame() *SnakeGame {
	// 获取浏览器DOM元素
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "snakeCanvas")
	ctx := canvas.Call("getContext", "2d")

	// 初始化蛇（画布中间，3节身体）
	startX := canvasWidth / 2 / gridSize * gridSize
	startY := canvasHeight / 2 / gridSize * gridSize
	snake := []Point{
		{X: startX, Y: startY},
		{X: startX - gridSize, Y: startY},
		{X: startX - 2*gridSize, Y: startY},
	}

	game := &SnakeGame{
		ctx:       ctx,
		canvas:    canvas,
		snake:     snake,
		direction: Right,
		nextDir:   Right,
		food:      Point{},
		score:     0,
		gameOver:  false,
	}

	// 生成初始食物
	game.spawnFood()
	// 注册键盘事件监听
	game.registerKeydownListener()

	return game
}

// 生成食物（随机位置，不与蛇身重叠）
func (g *SnakeGame) spawnFood() {
	rand.Seed(time.Now().UnixNano())
	for {
		x := rand.Intn(canvasWidth/gridSize) * gridSize
		y := rand.Intn(canvasHeight/gridSize) * gridSize
		food := Point{X: x, Y: y}

		// 检查是否与蛇身重叠
		overlap := false
		for _, seg := range g.snake {
			if seg == food {
				overlap = true
				break
			}
		}

		if !overlap {
			g.food = food
			break
		}
	}
}

// 注册键盘按下事件监听
func (g *SnakeGame) registerKeydownListener() {
	// 定义键盘事件处理函数
	keydownHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		key := e.Get("key").String()

		switch key {
		case "ArrowUp", "w", "W":
			if g.direction != Down {
				g.nextDir = Up
			}
		case "ArrowDown", "s", "S":
			if g.direction != Up {
				g.nextDir = Down
			}
		case "ArrowLeft", "a", "A":
			if g.direction != Right {
				g.nextDir = Left
			}
		case "ArrowRight", "d", "D":
			if g.direction != Left {
				g.nextDir = Right
			}
		case "r", "R":
			if g.gameOver {
				g.reset()
			}
		}
		return nil
	})

	// 绑定到window的keydown事件
	js.Global().Call("addEventListener", "keydown", keydownHandler)
	// 防止GC回收（Go Wasm中需保持引用）
	//defer keydownHandler.Release()
}

// 重置游戏
func (g *SnakeGame) reset() {
	startX := canvasWidth / 2 / gridSize * gridSize
	startY := canvasHeight / 2 / gridSize * gridSize
	g.snake = []Point{
		{X: startX, Y: startY},
		{X: startX - gridSize, Y: startY},
		{X: startX - 2*gridSize, Y: startY},
	}
	g.direction = Right
	g.nextDir = Right
	g.score = 0
	g.gameOver = false
	g.spawnFood()
}

// 更新游戏状态（移动蛇、检测碰撞、吃食物）
func (g *SnakeGame) update() {
	if g.gameOver {
		return
	}

	// 更新当前方向
	g.direction = g.nextDir

	// 获取蛇头，计算新蛇头
	head := g.snake[0]
	var newHead Point
	switch g.direction {
	case Up:
		newHead = Point{X: head.X, Y: head.Y - gridSize}
	case Down:
		newHead = Point{X: head.X, Y: head.Y + gridSize}
	case Left:
		newHead = Point{X: head.X - gridSize, Y: head.Y}
	case Right:
		newHead = Point{X: head.X + gridSize, Y: head.Y}
	}

	// 1. 检测撞墙（游戏结束）
	if newHead.X < 0 || newHead.X >= canvasWidth ||
		newHead.Y < 0 || newHead.Y >= canvasHeight {
		g.gameOver = true
		return
	}

	// 2. 检测撞到自己（游戏结束）
	for _, seg := range g.snake[1:] {
		if seg == newHead {
			g.gameOver = true
			return
		}
	}

	// 3. 添加新蛇头到头部
	g.snake = append([]Point{newHead}, g.snake...)

	// 4. 检测是否吃食物
	if newHead == g.food {
		g.score += 10
		g.spawnFood()
	} else {
		// 没吃到食物，删除蛇尾
		g.snake = g.snake[:len(g.snake)-1]
	}
}

// 渲染游戏画面（绘制蛇、食物、分数、游戏结束提示）
func (g *SnakeGame) render() {
	// 1. 清空画布
	g.ctx.Call("clearRect", 0, 0, canvasWidth, canvasHeight)

	// 2. 绘制蛇
	// 蛇头（绿色）
	g.ctx.Set("fillStyle", "#2ecc71")
	g.ctx.Call("fillRect", g.snake[0].X, g.snake[0].Y, gridSize-1, gridSize-1)

	// 蛇身（深绿色）
	g.ctx.Set("fillStyle", "#27ae60")
	for _, seg := range g.snake[1:] {
		g.ctx.Call("fillRect", seg.X, seg.Y, gridSize-1, gridSize-1)
	}

	// 3. 绘制食物（红色）
	g.ctx.Set("fillStyle", "#e74c3c")
	g.ctx.Call("fillRect", g.food.X, g.food.Y, gridSize-1, gridSize-1)

	// 4. 绘制分数
	g.ctx.Set("fillStyle", "#333333")
	g.ctx.Set("font", "20px Arial")
	g.ctx.Call("fillText", "Score: "+strconv.Itoa(g.score), 10, 30)

	// 5. 绘制游戏结束提示
	if g.gameOver {
		g.ctx.Set("fillStyle", "#ff0000")
		g.ctx.Set("font", "40px Arial")
		g.ctx.Call("fillText", "GAME OVER", canvasWidth/2-120, canvasHeight/2)
		g.ctx.Set("font", "20px Arial")
		g.ctx.Call("fillText", "Press R to Restart", canvasWidth/2-80, canvasHeight/2+40)
	}
}

// 游戏主循环
func (g *SnakeGame) gameLoop() {
	ticker := time.NewTicker(time.Millisecond * updateSpeed)
	defer ticker.Stop()

	for range ticker.C {
		g.update()
		g.render()
	}
}

// 辅助：Go没有内置strconv.Itoa的极简替代（避免额外导入）
func strconvItoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append(b, byte('0'+n%10))
		n /= 10
	}
	// 反转字节切片
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	// 初始化游戏
	game := NewSnakeGame()

	// 启动游戏循环（阻塞主线程，防止Wasm程序退出）
	game.gameLoop()

	// 保持程序运行（Go Wasm中main函数退出后程序终止）
	select {}
}
