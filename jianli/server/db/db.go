package db

import (
	"jianli-server/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.New(sqlite.Config{
		DSN: "file:jianli.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)",
		DriverName: "sqlite",
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.Template{}, &models.Resume{}, &models.Version{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedTemplates()

	log.Println("Database initialized successfully")
}

func seedTemplates() {
	var count int64
	DB.Model(&models.Template{}).Count(&count)
	if count > 0 {
		return
	}

	templates := []models.Template{
		{
			Name:        "简约商务模板",
			Description: "简洁大方的商务风格，适合求职各类岗位",
			Category:    "商务",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic_1", "type": "basic", "name": "基本信息", "visible": true, "order": 1, "removable": false},
					{"id": "edu_1", "type": "education", "name": "教育经历", "visible": true, "order": 2, "removable": true},
					{"id": "exp_1", "type": "experience", "name": "工作经历", "visible": true, "order": 3, "removable": true},
					{"id": "sk_1", "type": "skills", "name": "专业技能", "visible": true, "order": 4, "removable": true},
					{"id": "prj_1", "type": "projects", "name": "项目经验", "visible": true, "order": 5, "removable": true},
				},
				"basic": map[string]interface{}{
					"name": "张三", "phone": "138-0000-0000", "email": "zhangsan@example.com",
					"city": "北京", "job_intention": "高级后端开发工程师",
					"github": "https://github.com/zhangsan", "website": "https://zhangsan.dev",
				},
				"education": []map[string]interface{}{
					{"school": "清华大学", "degree": "硕士", "major": "计算机科学与技术",
						"start_date": "2018-09", "end_date": "2021-06",
						"description": "主修课程：数据结构、算法设计、分布式系统、机器学习。GPA: 3.8/4.0"},
					{"school": "北京大学", "degree": "本科", "major": "软件工程",
						"start_date": "2014-09", "end_date": "2018-06",
						"description": "连续三年获得校级一等奖学金，学生会技术部部长"},
				},
				"experience": []map[string]interface{}{
					{"company": "阿里巴巴集团", "position": "高级工程师",
						"start_date": "2023-03", "end_date": "至今",
						"description": "负责电商中台核心系统的设计与开发，支撑双十一亿级流量。",
						"highlights": []string{"主导订单系统重构，QPS提升300%", "优化数据库索引，查询耗时降低60%", "设计分布式事务方案，数据一致性达99.99%"}},
					{"company": "字节跳动", "position": "后端工程师",
						"start_date": "2021-07", "end_date": "2023-02",
						"description": "参与内容推荐系统的开发与维护，负责推荐召回模块。",
						"highlights": []string{"实现多通道召回策略，点击率提升15%", "开发实时特征计算平台，特征更新延迟从小时级降至秒级"}},
				},
				"skills": []map[string]interface{}{
					{"category": "后端开发", "items": "Go, Java, Python, gRPC, RESTful API"},
					{"category": "数据存储", "items": "MySQL, PostgreSQL, Redis, MongoDB, Elasticsearch"},
					{"category": "运维部署", "items": "Docker, Kubernetes, Nginx, CI/CD, Prometheus"},
				},
				"projects": []map[string]interface{}{
					{"name": "分布式电商订单系统", "role": "技术负责人",
						"start_date": "2023-06", "end_date": "2023-12",
						"description": "从零设计并实现支持日均千万级订单的分布式电商系统，包含交易、库存、支付等核心模块。",
						"highlights": []string{"采用微服务架构拆分6大领域服务", "引入Sentinel熔断降级，系统可用性达99.95%", "设计分库分表方案，支撑数据量超10亿"}},
					{"name": "实时用户画像平台", "role": "核心开发",
						"start_date": "2022-04", "end_date": "2022-10",
						"description": "基于Flink + Kafka构建实时用户行为分析平台，支持百万人级实时标签计算。",
						"highlights": []string{"优化Flink状态后端，内存占用降低40%", "设计实时-离线一体化标签计算架构"}},
				},
				"summary":    "",
				"evaluation": "",
			},
			StyleConfig: models.JSON{
				"primaryColor":  "#1a1a2e",
				"accentColor":   "#0f3460",
				"fontSize":      12,
				"lineSpacing":   1.5,
				"marginTop":     20,
				"marginBottom":  20,
				"marginLeft":    25,
				"marginRight":   25,
			},
		},
		{
			Name:        "创意设计模板",
			Description: "充满创意的设计风格，适合设计、创意类岗位",
			Category:    "创意",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic_1", "type": "basic", "name": "基本信息", "visible": true, "order": 1, "removable": false},
					{"id": "sum_1", "type": "summary", "name": "个人简介", "visible": true, "order": 2, "removable": true},
					{"id": "sk_1", "type": "skills", "name": "专业技能", "visible": true, "order": 3, "removable": true},
					{"id": "exp_1", "type": "experience", "name": "工作经历", "visible": true, "order": 4, "removable": true},
					{"id": "prj_1", "type": "projects", "name": "作品集", "visible": true, "order": 5, "removable": true},
					{"id": "edu_1", "type": "education", "name": "教育经历", "visible": true, "order": 6, "removable": true},
				},
				"basic": map[string]interface{}{
					"name": "李四", "phone": "139-0000-0000", "email": "lisi@example.com",
					"city": "上海", "job_intention": "资深UI设计师",
					"github": "", "website": "https://lisi.design",
				},
				"summary":    "7年UI设计经验，擅长移动端App设计和品牌视觉系统设计。曾主导过3款百万用户级产品的视觉设计，作品获得过红点设计奖和iF设计奖。对用户体验有深刻理解，善于将商业目标与用户需求结合。",
				"education": []map[string]interface{}{
					{"school": "中央美术学院", "degree": "本科", "major": "视觉传达设计",
						"start_date": "2014-09", "end_date": "2018-06",
						"description": "毕业作品入选年度优秀作品展"},
				},
				"experience": []map[string]interface{}{
					{"company": "腾讯", "position": "高级UI设计师",
						"start_date": "2021-01", "end_date": "至今",
						"description": "负责微信支付相关产品的视觉设计和设计规范建设。",
						"highlights": []string{"主导微信支付9.0版本视觉改版", "搭建支付设计系统，覆盖200+组件", "设计的刷脸支付界面获2022红点设计奖"}},
					{"company": "美团", "position": "UI设计师",
						"start_date": "2018-07", "end_date": "2020-12",
						"description": "负责美团外卖App用户端视觉设计。",
						"highlights": []string{"完成外卖V8版本视觉整体升级", "设计外卖品牌IP形象，日活提升5%"}},
				},
				"skills": []map[string]interface{}{
					{"category": "设计软件", "items": "Figma, Sketch, Adobe XD, Photoshop, Illustrator"},
					{"category": "设计能力", "items": "界面设计, 交互设计, 设计系统, 品牌设计, 插画"},
					{"category": "其他技能", "items": "用户研究, 动效设计(Principle), HTML/CSS基础"},
				},
				"projects": []map[string]interface{}{
					{"name": "微信支付9.0视觉改版", "role": "主设计师",
						"start_date": "2022-03", "end_date": "2022-09",
						"description": "作为主设计师主导微信支付9.0版本的整体视觉升级，提升品牌年轻感和科技感。",
						"highlights": []string{"用户满意度从82分提升至89分", "界面NPS值提升12个点"}},
					{"name": "企业设计系统建设", "role": "项目负责人",
						"start_date": "2023-01", "end_date": "2023-06",
						"description": "从0到1搭建跨业务线的企业级设计系统，包含设计规范、组件库、文档站。",
						"highlights": []string{"设计效率提升50%", "视觉一致性评分从71%提升至94%"}},
				},
				"evaluation": "",
			},
			StyleConfig: models.JSON{
				"primaryColor":  "#2c3e50",
				"accentColor":   "#e74c3c",
				"fontSize":      12,
				"lineSpacing":   1.6,
				"marginTop":     15,
				"marginBottom":  15,
				"marginLeft":    20,
				"marginRight":   20,
			},
		},
		{
			Name:        "技术工程师模板",
			Description: "专业技术风格，突出技能和项目经验",
			Category:    "技术",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic_1", "type": "basic", "name": "基本信息", "visible": true, "order": 1, "removable": false},
					{"id": "sum_1", "type": "summary", "name": "技术简介", "visible": true, "order": 2, "removable": true},
					{"id": "sk_1", "type": "skills", "name": "技术栈", "visible": true, "order": 3, "removable": true},
					{"id": "prj_1", "type": "projects", "name": "项目经验", "visible": true, "order": 4, "removable": true},
					{"id": "exp_1", "type": "experience", "name": "工作经历", "visible": true, "order": 5, "removable": true},
					{"id": "edu_1", "type": "education", "name": "教育经历", "visible": true, "order": 6, "removable": true},
				},
				"basic": map[string]interface{}{
					"name": "王五", "phone": "137-0000-0000", "email": "wangwu@example.com",
					"city": "深圳", "job_intention": "技术专家 / 架构师",
					"github": "https://github.com/wangwu", "website": "",
				},
				"summary":    "10年+软件开发经验，5年技术团队管理经验。精通高并发、高可用系统架构设计，主导过多个日活千万级系统的架构设计与技术选型。技术栈全面，前后端兼修，对云原生、微服务、大数据有丰富实战经验。",
				"education": []map[string]interface{}{
					{"school": "上海交通大学", "degree": "硕士", "major": "计算机应用技术",
						"start_date": "2011-09", "end_date": "2014-03",
						"description": "研究方向：分布式计算与云计算，发表SCI论文2篇"},
				},
				"experience": []map[string]interface{}{
					{"company": "华为", "position": "技术专家",
						"start_date": "2019-05", "end_date": "至今",
						"description": "负责华为云PaaS平台核心服务的架构设计与技术攻坚。",
						"highlights": []string{"主导Service Mesh落地，服务间调用P99降低40%", "设计Serverless函数计算平台，冷启动时间<200ms", "培养技术团队15人，输出专家3人"}},
					{"company": "京东", "position": "高级架构师",
						"start_date": "2016-04", "end_date": "2019-04",
						"description": "负责京东商城商品中心架构设计与优化。",
						"highlights": []string{"完成商品详情页缓存架构重构，QPS从2万提升至15万", "推进全链路压测体系建设，系统承载能力翻倍"}},
				},
				"skills": []map[string]interface{}{
					{"category": "编程语言", "items": "Java, Go, C++, Python, TypeScript"},
					{"category": "架构设计", "items": "微服务, DDD领域驱动, 分布式事务, 高可用设计, 性能优化"},
					{"category": "云原生", "items": "Kubernetes, Docker, Istio, Serverless, Operator开发"},
					{"category": "中间件", "items": "Nacos, Sentinel, RocketMQ, Kafka, Elasticsearch, Redis Cluster"},
				},
				"projects": []map[string]interface{}{
					{"name": "下一代微服务治理平台", "role": "架构师",
						"start_date": "2022-01", "end_date": "2023-06",
						"description": "基于Istio构建的全链路微服务治理平台，覆盖5000+微服务，支持灰度发布、熔断降级、流量镜像、故障注入等能力。",
						"highlights": []string{"基于eBPF实现Sidecar资源占用降低50%", "自研灰度发布引擎，上线风险降低80%"}},
					{"name": "实时数仓建设", "role": "技术负责人",
						"start_date": "2020-03", "end_date": "2021-08",
						"description": "基于Flink + Kafka + Doris构建公司级实时数仓，支撑T+0业务报表和实时大屏。",
						"highlights": []string{"日处理数据量超100TB", "端到端延迟从小时级降至秒级"}},
				},
				"evaluation": "",
			},
			StyleConfig: models.JSON{
				"primaryColor":  "#000000",
				"accentColor":   "#0066cc",
				"fontSize":      11,
				"lineSpacing":   1.4,
				"marginTop":     20,
				"marginBottom":  20,
				"marginLeft":    25,
				"marginRight":   25,
			},
		},
	}

	for _, t := range templates {
		DB.Create(&t)
	}

	log.Println("Template data seeded")
}
