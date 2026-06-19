package main

import (
	"fmt"
	"log"
	"shuiti/database"
	"shuiti/models"
	"shuiti/router"
)

func seedData() {
	var count int64
	database.DB.Model(&models.Question{}).Count(&count)
	if count > 0 {
		return
	}

	questions := []models.Question{
		{
			Type:       models.SingleChoice,
			Content:    "Go语言中，以下哪个是正确声明变量的方式？",
			Options:    "[\"A. var x int = 10\",\"B. variable x = 10\",\"C. int x = 10\",\"D. x :=: 10\"]",
			Answer:     "A",
			Analysis:   "Go语言使用var关键字声明变量，或使用:=短变量声明。正确格式：var x int = 10 或 x := 10",
			Difficulty: 1,
			Category:   "Go基础",
		},
		{
			Type:       models.SingleChoice,
			Content:    "Gin框架的默认端口是？",
			Options:    "[\"A. 3000\",\"B. 8000\",\"C. 8080\",\"D. 5000\"]",
			Answer:     "C",
			Analysis:   "Gin框架默认使用8080端口启动HTTP服务",
			Difficulty: 1,
			Category:   "Gin框架",
		},
		{
			Type:       models.SingleChoice,
			Content:    "下列哪个不是Go语言的基本数据类型？",
			Options:    "[\"A. int\",\"B. string\",\"C. float\",\"D. boolean\"]",
			Answer:     "D",
			Analysis:   "Go语言中的布尔类型是bool，而不是boolean",
			Difficulty: 1,
			Category:   "Go基础",
		},
		{
			Type:       models.MultipleChoice,
			Content:    "以下哪些是Go语言的循环语句？（多选）",
			Options:    "[\"A. for\",\"B. while\",\"C. range\",\"D. foreach\"]",
			Answer:     "AC",
			Analysis:   "Go语言只有for循环，但可以模拟while的用法。range用于遍历数组、切片、映射等。Go没有foreach关键字",
			Difficulty: 2,
			Category:   "Go基础",
		},
		{
			Type:       models.MultipleChoice,
			Content:    "Gin框架支持以下哪些HTTP方法？（多选）",
			Options:    "[\"A. GET\",\"B. POST\",\"C. PUT\",\"D. PATCH\"]",
			Answer:     "ABCD",
			Analysis:   "Gin框架支持所有标准HTTP方法：GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS等",
			Difficulty: 1,
			Category:   "Gin框架",
		},
		{
			Type:       models.TrueFalse,
			Content:    "Go语言中，切片（slice）是值类型。",
			Options:    "[\"A. 正确\",\"B. 错误\"]",
			Answer:     "B",
			Analysis:   "切片是引用类型，它包含指向底层数组的指针、长度和容量。传递切片时传递的是引用",
			Difficulty: 2,
			Category:   "Go基础",
		},
		{
			Type:       models.TrueFalse,
			Content:    "SQLite是一种轻量级的文件型数据库，不需要独立的服务进程。",
			Options:    "[\"A. 正确\",\"B. 错误\"]",
			Answer:     "A",
			Analysis:   "SQLite是嵌入式数据库，以单个文件存储，无需服务器进程，适合轻量级应用",
			Difficulty: 1,
			Category:   "数据库",
		},
		{
			Type:       models.SingleChoice,
			Content:    "Vue3中，用于创建响应式对象的函数是？",
			Options:    "[\"A. reactive()\",\"B. ref()\",\"C. computed()\",\"D. watch()\"]",
			Answer:     "A",
			Analysis:   "reactive()用于创建响应式对象，ref()用于创建响应式基本类型，computed()用于计算属性",
			Difficulty: 1,
			Category:   "Vue3",
		},
		{
			Type:       models.SingleChoice,
			Content:    "以下哪个不是Vue3的生命周期钩子？",
			Options:    "[\"A. onMounted\",\"B. onCreated\",\"C. onUnmounted\",\"D. onUpdated\"]",
			Answer:     "B",
			Analysis:   "Vue3组合式API中没有onCreated，setup函数本身相当于created阶段。生命周期包括：onMounted、onUpdated、onUnmounted等",
			Difficulty: 2,
			Category:   "Vue3",
		},
		{
			Type:       models.MultipleChoice,
			Content:    "以下哪些是GORM支持的关联关系？（多选）",
			Options:    "[\"A. HasOne\",\"B. HasMany\",\"C. BelongsTo\",\"D. ManyToMany\"]",
			Answer:     "ABCD",
			Analysis:   "GORM支持所有常见的数据库关联：一对一(HasOne/BelongsTo)、一对多(HasMany)、多对多(ManyToMany)",
			Difficulty: 2,
			Category:   "数据库",
		},
		{
			Type:       models.TrueFalse,
			Content:    "Vue3中的Composition API只能在<script setup>中使用。",
			Options:    "[\"A. 正确\",\"B. 错误\"]",
			Answer:     "B",
			Analysis:   "Composition API可以在setup()函数中使用，也可以在<script setup>语法糖中使用，后者更简洁",
			Difficulty: 2,
			Category:   "Vue3",
		},
		{
			Type:       models.SingleChoice,
			Content:    "Gin框架中，用于获取JSON请求体并绑定到结构体的方法是？",
			Options:    "[\"A. c.Bind()\",\"B. c.ShouldBindJSON()\",\"C. c.JSON()\",\"D. c.Query()\"]",
			Answer:     "B",
			Analysis:   "ShouldBindJSON()用于将JSON请求体绑定到结构体，Bind()会根据Content-Type自动选择，JSON()用于返回JSON响应",
			Difficulty: 1,
			Category:   "Gin框架",
		},
		{
			Type:       models.SingleChoice,
			Content:    "Go语言中，defer语句的执行顺序是？",
			Options:    "[\"A. 先进先出(FIFO)\",\"B. 后进先出(LIFO)\",\"C. 随机执行\",\"D. 按优先级执行\"]",
			Answer:     "B",
			Analysis:   "defer语句遵循后进先出(LIFO)原则，最后声明的defer最先执行，类似栈结构",
			Difficulty: 2,
			Category:   "Go基础",
		},
		{
			Type:       models.MultipleChoice,
			Content:    "以下哪些是HTTP的幂等方法？（多选）",
			Options:    "[\"A. GET\",\"B. POST\",\"C. PUT\",\"D. DELETE\"]",
			Answer:     "ACD",
			Analysis:   "GET、PUT、DELETE是幂等方法，多次执行结果相同。POST不是幂等的，每次执行可能创建新资源",
			Difficulty: 2,
			Category:   "网络基础",
		},
		{
			Type:       models.TrueFalse,
			Content:    "在Go语言中，channel可以用于goroutine之间的通信。",
			Options:    "[\"A. 正确\",\"B. 错误\"]",
			Answer:     "A",
			Analysis:   "Channel是Go语言中goroutine之间通信的核心机制，遵循\"不要通过共享内存通信，通过通信共享内存\"的设计哲学",
			Difficulty: 1,
			Category:   "Go基础",
		},
	}

	database.DB.Create(&questions)
	log.Printf("Seeded %d questions successfully", len(questions))
}

func main() {
	database.InitDB()
	seedData()

	r := router.SetupRouter()

	fmt.Println("========================================")
	fmt.Println("  刷题系统服务启动")
	fmt.Println("  服务地址: http://localhost:8080")
	fmt.Println("  API前缀:  http://localhost:8080/api")
	fmt.Println("========================================")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
